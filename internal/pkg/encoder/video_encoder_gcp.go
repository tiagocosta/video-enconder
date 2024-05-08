package encoder

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/tiagocosta/video-enconder/configs"
)

type VideoEncoderGCP struct{}

func (videoEncoder *VideoEncoderGCP) Download(filePath string, videoID string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(configs.Config().BucketName)
	obj := bkt.Object(filePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(configs.Config().LocalStoragePath + "/" + videoID + ".mp4")

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has been stored", videoID)

	return nil
}

func (videoEncoder *VideoEncoderGCP) Fragment(videoID string) error {
	localStoragePath := configs.Config().LocalStoragePath

	err := os.Mkdir(localStoragePath+"/"+videoID, os.ModePerm)
	if err != nil {
		return err
	}

	source := localStoragePath + "/" + videoID + ".mp4"
	target := localStoragePath + "/" + videoID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (videoEncoder *VideoEncoderGCP) Encode(videoID string) error {
	localStoragePath := configs.Config().LocalStoragePath

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, localStoragePath+"/"+videoID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, localStoragePath+"/"+videoID)
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

func (videoEncoder *VideoEncoderGCP) Upload(videoID string) error {
	videoPath := configs.Config().LocalStoragePath + "/" + videoID
	paths, err := loadPaths(videoPath)
	if err != nil {
		return err
	}

	in := make(chan string, runtime.NumCPU())
	out := make(chan string)
	concurrency, _ := strconv.Atoi(configs.Config().ConcurrencyUpload)

	uploadClient, ctx, err := newUploadClient()
	if err != nil {
		return err
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
			return errors.New(path + " " + result)
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

func (videoEncoder *VideoEncoderGCP) CleanupFiles(videoID string) error {
	localStoragePath := configs.Config().LocalStoragePath

	err := os.Remove(localStoragePath + "/" + videoID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", videoID, ".mp4")
		return err
	}

	err = os.Remove(localStoragePath + "/" + videoID + ".frag")
	if err != nil {
		log.Println("error removing frag ", videoID, ".frag")
		return err
	}

	err = os.RemoveAll(localStoragePath + "/" + videoID)
	if err != nil {
		log.Println("error removing mp4 ", videoID, ".mp4")
		return err
	}

	log.Println("files have been removed: ", videoID)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}
