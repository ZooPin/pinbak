package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pinbak/cmd/pinbak/helper"
)

var restoreCmd = &cobra.Command{
	Use:   "restore [repository-name]",
	Short: "Restore all file in repository.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one repository")
		}

		return nil
	},
	Run: restoreFunc,
}

var restoreAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Restore all file in repository.",
	Run:   restoreAllFunc,
}

func init() {
	restoreCmd.AddCommand(restoreAllCmd)
	rootCmd.AddCommand(restoreCmd)
}

func restoreFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Add error: ", err)
	}
	for _, repo := range args {
		if !mover.Config.CheckRepository(repo) {
			log.Print("Repository: ", repo, " does not exist.")
			continue
		}
		err = mover.Restore(repo)
		if err != nil {
			log.Print("Error restore ", repo, ": ", err)
			continue
		}
		fmt.Println("Repository ", repo, " restored.")
	}
}

func restoreAllFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Add error: ", err)
	}
	for repo, _ := range mover.Config.Repository {
		if !mover.Config.CheckRepository(repo) {
			log.Print("Repository: ", repo, " does not exist.")
			continue
		}
		err = mover.Restore(repo)
		if err != nil {
			log.Print("Error restore ", repo, ": ", err)
			continue
		}
		fmt.Println("Repository ", repo, " restored.")
	}
}
