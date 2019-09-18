package helper

import (
	"fmt"
	"os"
	"pinbak"
)

func GetConfig() (pinbak.Config, error) {
	return pinbak.LoadConfig(PinbakPath())
}

func GetGitHelper() (pinbak.GitHelper, error) {
	config, err := GetConfig()
	if err != nil {
		return pinbak.GitHelper{}, err
	}
	return pinbak.CreateGit(config), nil
}

func GetMover() (pinbak.Mover, error) {
	git, err := GetGitHelper()
	if err != nil {
		return pinbak.Mover{}, err
	}
	return pinbak.CreateMover(git.Config, git), nil
}

func homeDir() string {
	s, _ := os.UserHomeDir()
	return s
}

func PinbakPath() string {
	return fmt.Sprint(homeDir(), "/.pinbak")
}
