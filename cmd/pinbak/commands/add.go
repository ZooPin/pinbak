package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pinbak/cmd/pinbak/helper"
)

var addCmd = &cobra.Command{
	Use:   "add [repository] [file] [file...]",
	Short: "Add a file to backup.",
	Args:  cobra.MinimumNArgs(2),
	Run:   addFunc,
}

var addRepoCmd = &cobra.Command{
	Use:   "repository [name] [url]",
	Short: "Add a repository to pinbak.",
	Args:  cobra.MinimumNArgs(2),
	Run:   addRepoFunc,
}

func init() {
	addCmd.AddCommand(addRepoCmd)
	rootCmd.AddCommand(addCmd)
}

func addFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Add error: ", err)
	}

	for i := 1; i < len(args); i++ {
		err = mover.Add(args[i], args[0])
		if err != nil {
			log.Print("Add error in ", args[0], " file ", args[i])
		}
	}
	err = mover.Git.CommitAndPush(args[0])
	if err != nil {
		log.Print("Add error: ", err)
	}

	fmt.Println("Done.")
}

func addRepoFunc(cmd *cobra.Command, args []string) {
	git, err := helper.GetGitHelper()
	if err != nil {
		log.Fatal("Add repo error:", err)
	}
	err = git.Clone(args[0], args[1])
	if err != nil {
		log.Fatal("Add repo  error:", err)
	}
	fmt.Println("Repository", args[0], "added.")
}
