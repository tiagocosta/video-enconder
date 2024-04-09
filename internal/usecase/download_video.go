package usecase

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

type DownloadVideoInputDTO struct {
	BucketName string
	FilePath   string
	VideoID    string
}

type DownloadVideoUseCase struct {
}

func NewDownloadVideoUseCase() *DownloadVideoUseCase {
	return &DownloadVideoUseCase{}
}

func (c *DownloadVideoUseCase) Execute(input DownloadVideoInputDTO) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(input.BucketName)
	obj := bkt.Object(input.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("LOCAL_STORAGE_PATH") + "/" + input.VideoID + ".mp4")

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has been stored", input.VideoID)

	return nil
}
