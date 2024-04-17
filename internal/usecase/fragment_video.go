package usecase

import (
	"os"
	"os/exec"
)

type FragmentVideoInputDTO struct {
	VideoID string
}

type FragmentVideoUseCase struct{}

func NewFragmentVideoUseCase() *FragmentVideoUseCase {
	return &FragmentVideoUseCase{}
}

func (c *FragmentVideoUseCase) Execute(input FragmentVideoInputDTO) error {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")

	err := os.Mkdir(localStoragePath+"/"+input.VideoID, os.ModePerm)
	if err != nil {
		return err
	}

	source := localStoragePath + "/" + input.VideoID + ".mp4"
	target := localStoragePath + "/" + input.VideoID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}
