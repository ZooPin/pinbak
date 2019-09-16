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

}
