package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pinbak/cmd/pinbak/helper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all file in repository.",
	Run:   listFunc,
}

var listRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "List all file in all repositories.",
	Run:   listRepository,
}

var listFileCmd = &cobra.Command{
	Use:   "file [...repo]",
	Short: "List all file in a repository.",
	Args:  cobra.MinimumNArgs(1),
	Run:   listFiles,
}

func init() {
	listCmd.AddCommand(listFileCmd)
	listCmd.AddCommand(listRepoCmd)
	rootCmd.AddCommand(listCmd)
}

func listFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Error list:", err)
	}
	for repo, url := range mover.Config.Repository {
		fmt.Println(repo, ":", url)
		files, err := mover.List(repo)
		if err != nil {
			log.Println("Error list ", repo, ": ", err)
			continue
		}
		for id, restorePath := range files {
			fmt.Println("    - ", id, " : ", restorePath)
		}
	}
}

func listRepository(cmd *cobra.Command, args []string) {
	config, err := helper.GetConfig()
	if err != nil {
		log.Fatal("Error list:", err)
	}
	for repo, url := range config.Repository {
		fmt.Println(repo, " : ", url)
	}
}

func listFiles(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Error list:", err)
	}
	for _, repo := range args {
		if !mover.Config.CheckRepository(repo) {
			log.Println("Repositoy", repo, "not found.")
			continue
		}
		files, err := mover.List(repo)
		if err != nil {
			log.Println("Error list ", repo, ": ", err)
			continue
		}
		fmt.Println(repo, " : ", mover.Config.Repository[repo])
		for id, restorePath := range files {
			fmt.Println("    - ", id, " : ", restorePath)
		}
	}
}
