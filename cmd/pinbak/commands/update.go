package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pinbak/cmd/pinbak/helper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all files in all repositories.",
	Run:   updateFunc,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Add error: ", err)
	}
	for repo, _ := range mover.Config.Repository {
		err = mover.Update(repo)
		if err != nil {
			log.Println("Update error: ", err)
		}
		err = mover.Git.CommitAndPush(repo)
		if err != nil {
			log.Println("Update error: ", err)
		}
	}
	fmt.Println("Done.")
}
