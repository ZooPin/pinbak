package main

import (
	"fmt"
	"log"
	"os"
	"pinbak"
)

func main() {
	homePath, _ := os.UserHomeDir()
	path := fmt.Sprint(homePath, "/.pinbak")
	configPath := fmt.Sprint(path, "/config")
	var config pinbak.Config
	config, err := config.Load(configPath)
	if err != nil {
		log.Fatal("Error: ", config)
	}

	git := pinbak.GitHelper{
		Path:   path,
		Config: &config,
	}

	mover := pinbak.CreateMover(config, git, path)

	err = mover.Add("/home/pingouin/.zshrc", "pinbak-test")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	err = mover.Add("/home/pingouin/git/.config", "pinbak-test")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	err = git.CommitAndPush("pinbak-test")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Print("Done")

}
