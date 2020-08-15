package commands

import (
	"fmt"
	"github.com/pngouin/pinbak/helper"
	"github.com/spf13/cobra"
	"log"
)

var removeCmd = &cobra.Command{
	Use:   "remove [id] [id...]",
	Short: "Remove items or repository from the backup.",
	Args:  cobra.MinimumNArgs(1),
	Run:   removeFunc,
}

var removeRepoCmd = &cobra.Command{
	Use:   "repo [name]",
	Short: "Remove repository from the backup.",
	Args:  cobra.MaximumNArgs(1),
	Run:   removeRepoFunc,
}

func init() {
	removeCmd.AddCommand(removeRepoCmd)
	rootCmd.AddCommand(removeCmd)
}

func removeRepoFunc(cmd *cobra.Command, args []string) {
	config, err := helper.GetConfig()
	if err != nil {
		log.Fatal("Remove error: ", err)
	}
	err = config.RemoveRepository(args[0])
	if err != nil {
		log.Println("Remove error: ", err)
	}
	fmt.Println("Repository", args[0], "removed.")
}

func removeFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Remove error: ", err)
	}

	var repos []string
	for i := 0; i < len(args); i++ {
		r, err := mover.Remove(args[i])
		if err != nil {
			log.Print("Remove error with items: ", args[i], ": ", err)
			continue
		}
		repos = append(repos, r)
	}

	for _, repo := range repos {
		err = mover.Git.CommitAndPush(repo)
		if err != nil {
			log.Fatalln("Remove error: ", err)
		}
	}
	fmt.Println("Done.")
}
