package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
)

type ExecuteJobInputDTO struct {
	FilePath   string
	ResourceID string
}

type ExecuteJobUseCase struct {
	VideoRepository *database.VideoRepository
	JobRepository   *database.JobRepository
	JobCompleted    *event.JobCompleted
	EventDispatcher *events.EventDispatcher
}

func NewExecuteJobUseCase(
	videoRepository *database.VideoRepository,
	jobRepository *database.JobRepository,
	jobCompleted *event.JobCompleted,
	eventDispatcher *events.EventDispatcher,
) *ExecuteJobUseCase {
	return &ExecuteJobUseCase{
		VideoRepository: videoRepository,
		JobRepository:   jobRepository,
		JobCompleted:    jobCompleted,
		EventDispatcher: eventDispatcher,
	}
}

func (uc *ExecuteJobUseCase) Execute(input ExecuteJobInputDTO) error {
	video, _ := entity.NewVideo(uuid.NewString(), input.ResourceID, input.FilePath)
	uc.VideoRepository.Save(video)

	bucketName := configs.Config().BucketName
	job, _ := entity.NewJob(bucketName, entity.Pending, video)
	uc.JobRepository.Save(job)

	// workers, _ := strconv.Atoi(configs.Config().ConcurrencyWorkers)
	// for workerId := 0; workerId < workers; workerId++ {
	// uc.worker(1, job)
	// }

	err := uc.download(job)
	if err != nil {
		return err
	}

	err = uc.fragment(job)
	if err != nil {
		return err
	}

	err = uc.encode(job)
	if err != nil {
		return err
	}

	err = uc.upload(job)
	if err != nil {
		return err
	}

	err = uc.cleanupFiles(job)
	if err != nil {
		return err
	}

	err = uc.complete(job)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ExecuteJobUseCase) updateJob(job *entity.Job) error {
	err := uc.JobRepository.Update(job)

	if err != nil {
		return err
	}

	return nil
}

func (uc *ExecuteJobUseCase) failJob(job *entity.Job, error error) error {
	job.Fail()
	err := uc.updateJob(job)

	if err != nil {
		return err
	}

	return error
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}

func (uc *ExecuteJobUseCase) download(job *entity.Job) error {
	job.StartVideoDownload()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}
	fmt.Println(job)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return uc.failJob(job, err)
	}

	bkt := client.Bucket(configs.Config().BucketName)
	obj := bkt.Object(job.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return uc.failJob(job, err)
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return uc.failJob(job, err)
	}

	f, err := os.Create(configs.Config().LocalStoragePath + "/" + job.VideoID + ".mp4")

	if err != nil {
		return uc.failJob(job, err)
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return uc.failJob(job, err)
	}

	log.Printf("video %v has been stored", job.VideoID)

	return nil
}

func (uc *ExecuteJobUseCase) fragment(job *entity.Job) error {
	job.StartVideoFragmentation()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}

	localStoragePath := configs.Config().LocalStoragePath

	err = os.Mkdir(localStoragePath+"/"+job.VideoID, os.ModePerm)
	if err != nil {
		return uc.failJob(job, err)
	}

	source := localStoragePath + "/" + job.VideoID + ".mp4"
	target := localStoragePath + "/" + job.VideoID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (uc *ExecuteJobUseCase) encode(job *entity.Job) error {
	job.StartVideoEncoding()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}

	localStoragePath := configs.Config().LocalStoragePath

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, localStoragePath+"/"+job.VideoID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, localStoragePath+"/"+job.VideoID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (uc *ExecuteJobUseCase) upload(job *entity.Job) error {
	job.StartVideoUpload()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}

	videoPath := configs.Config().LocalStoragePath + "/" + job.VideoID
	paths, err := loadPaths(videoPath)
	if err != nil {
		return uc.failJob(job, err)
	}

	in := make(chan string, runtime.NumCPU())
	out := make(chan string)
	concurrency, _ := strconv.Atoi(configs.Config().ConcurrencyUpload)

	uploadClient, ctx, err := newUploadClient()
	if err != nil {
		return uc.failJob(job, err)
	}

	for process := 0; process < concurrency; process++ {
		go uploadWorker(ctx, uploadClient, in, out)
	}

	go func() {
		for _, path := range paths {
			in <- path
		}
	}()

	for _, path := range paths {
		result := <-out
		if result != "" {
			return uc.failJob(job, errors.New(path+" "+result))
		}
	}

	return nil
}

func loadPaths(videoPath string) ([]string, error) {
	var paths []string
	err := filepath.Walk(videoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func newUploadClient() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func uploadWorker(ctx context.Context, upClient *storage.Client, in chan string, out chan string) {
	for path := range in {
		err := uploadObject(ctx, upClient, path)

		if err != nil {
			log.Printf("error while uploading: %v. Error: %v", path, err)
			out <- err.Error()
		}

		out <- ""
	}
}

func uploadObject(ctx context.Context, upClient *storage.Client, objectPath string) error {
	localStoragePath := configs.Config().LocalStoragePath
	outputBucket := configs.Config().BucketName

	path := strings.Split(objectPath, localStoragePath+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := upClient.Bucket(outputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

func (uc *ExecuteJobUseCase) cleanupFiles(job *entity.Job) error {
	job.CleanupVideoFiles()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}

	localStoragePath := configs.Config().LocalStoragePath

	err = os.Remove(localStoragePath + "/" + job.VideoID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", job.VideoID, ".mp4")
		return uc.failJob(job, err)
	}

	err = os.Remove(localStoragePath + "/" + job.VideoID + ".frag")
	if err != nil {
		log.Println("error removing frag ", job.VideoID, ".frag")
		return uc.failJob(job, err)
	}

	err = os.RemoveAll(localStoragePath + "/" + job.VideoID)
	if err != nil {
		log.Println("error removing mp4 ", job.VideoID, ".mp4")
		return uc.failJob(job, err)
	}

	log.Println("files have been removed: ", job.VideoID)

	return nil
}

func (uc *ExecuteJobUseCase) complete(job *entity.Job) error {
	job.Complete()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}

	return nil
}
