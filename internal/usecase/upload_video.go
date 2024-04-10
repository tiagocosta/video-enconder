package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type UploadVideoInputDTO struct {
	BucketName  string
	VideoPath   string
	Concurrency int
	DoneUpload  chan string
}

type UploadVideoUseCase struct {
	Paths        []string
	UploadClient *storage.Client
	Ctx          context.Context
}

func NewUploadVideoUseCase() (*UploadVideoUseCase, error) {
	uploadClient, ctx, err := newUploadClient()
	if err != nil {
		return nil, err
	}
	return &UploadVideoUseCase{
		UploadClient: uploadClient,
		Ctx:          ctx,
	}, nil
}

func (uc *UploadVideoUseCase) Execute(input UploadVideoInputDTO) error {
	in := make(chan int, runtime.NumCPU())
	out := make(chan string)

	err := uc.loadPaths(input.VideoPath)
	if err != nil {
		return err
	}

	for process := 0; process < input.Concurrency; process++ {
		go uc.uploadWorker(in, out, input)
	}

	go func() {
		for x := 0; x < len(uc.Paths); x++ {
			in <- x
		}
	}()

	countDoneWorker := 0
	for r := range out {
		countDoneWorker++

		if r != "" {
			input.DoneUpload <- r
			break
		}

		if countDoneWorker == len(uc.Paths) {
			close(in)
		}
	}

	return nil
}

func (uc *UploadVideoUseCase) loadPaths(videoPath string) error {
	var paths []string
	err := filepath.Walk(videoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err, "AAAAAAAAAAAAAAAAAAAAA")
			return nil
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	uc.Paths = paths
	fmt.Println(uc.Paths)

	return nil
}

func newUploadClient() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func (uc *UploadVideoUseCase) uploadWorker(in chan int, out chan string, input UploadVideoInputDTO) {

	for x := range in {
		err := uc.uploadObject(uc.Paths[x], input.BucketName)

		if err != nil {
			log.Printf("error while uploading: %v. Error: %v", uc.Paths[x], err)
			out <- err.Error()
		}

		out <- ""
	}

	out <- "upload completed"
}

func (uc *UploadVideoUseCase) uploadObject(objectPath string, outputBucket string) error {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")

	path := strings.Split(objectPath, localStoragePath+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := uc.UploadClient.Bucket(outputBucket).Object(path[1]).NewWriter(uc.Ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}
