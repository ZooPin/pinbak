package helper

import (
	"os"
	"path"

	"github.com/ZooPin/pinbak/manager"
)

func GetConfig() (manager.Config, error) {
	return manager.LoadConfig(PinbakPath())
}

func GetGitHelper() (manager.GitHelper, error) {
	config, err := GetConfig()
	if err != nil {
		return manager.GitHelper{}, err
	}
	return manager.CreateGit(config), nil
}

func GetMover() (manager.Mover, error) {
	git, err := GetGitHelper()
	if err != nil {
		return manager.Mover{}, err
	}
	return manager.CreateMover(git.Config, git), nil
}

func homeDir() string {
	s, _ := os.UserHomeDir()
	return s
}

func PinbakPath() string {
	return path.Join(homeDir(), ".pinbak")
}
