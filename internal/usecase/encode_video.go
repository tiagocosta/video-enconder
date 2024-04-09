package usecase

import (
	"os"
	"os/exec"
)

type EncodeVideoInputDTO struct {
	VideoID string
}

type EncodeVideoUseCase struct{}

func NewEncodeVideoUseCase() *EncodeVideoUseCase {
	return &EncodeVideoUseCase{}
}

func (c *EncodeVideoUseCase) Execute(input EncodeVideoInputDTO) error {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, localStoragePath+"/"+input.VideoID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, localStoragePath+"/"+input.VideoID)
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
