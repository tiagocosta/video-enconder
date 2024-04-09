package usecase

import (
	"log"
	"os"
)

type CleanVideoInputDTO struct {
	VideoID string
}

type CleanVideoUseCase struct{}

func NewCleanVideoUseCase() *CleanVideoUseCase {
	return &CleanVideoUseCase{}
}

func (c *CleanVideoUseCase) Execute(input CleanVideoInputDTO) error {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")

	err := os.Remove(localStoragePath + "/" + input.VideoID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", input.VideoID, ".mp4")
		return err
	}

	err = os.Remove(localStoragePath + "/" + input.VideoID + ".frag")
	if err != nil {
		log.Println("error removing frag ", input.VideoID, ".frag")
		return err
	}

	err = os.RemoveAll(localStoragePath + "/" + input.VideoID)
	if err != nil {
		log.Println("error removing mp4 ", input.VideoID, ".mp4")
		return err
	}

	log.Println("files have been removed: ", input.VideoID)

	return nil
}
