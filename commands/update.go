package commands

import (
	"fmt"
	"github.com/pngouin/pinbak/helper"
	"github.com/spf13/cobra"
	"log"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all items in all repositories.",
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
	for repo := range mover.Config.Repository {
		errs := mover.Update(repo)
		if len(errs) > 0 {
			for _, err := range errs {
				log.Println("Update error with repo ", repo, ":", err)
			}
		}
		err = mover.Git.CommitAndPush(repo)
		if err != nil {
			log.Println("Update error with repo ", repo, ":", err)
		}
	}
	fmt.Println("Done.")
}
