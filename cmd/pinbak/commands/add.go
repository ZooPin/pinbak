package commands

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"pinbak"
)

var addCmd = &cobra.Command{
	Use:   "add [repository-name] [file to add]",
	Short: "Add a file to backup.",
	Run:   addFunc,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addFunc(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Repository name and file need to be provided.")
	}

	config, err := pinbak.LoadConfig(configPath())
	if err != nil {
		log.Fatal("Add error: ", err)
	}
	git := pinbak.CreateGit(config)
	mover := pinbak.CreateMover(config, git)

	err = mover.Add(args[0], args[1])
	if err != nil {
		return err
	}

}
