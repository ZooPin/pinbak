package commands

import (
	"fmt"
	"github.com/ZooPin/pinbak/helper"
	"github.com/spf13/cobra"
	"log"
)

var removeCmd = &cobra.Command{
	Use:   "remove [id] [id...]",
	Short: "Remove items from the backup.",
	Args:  cobra.MinimumNArgs(1),
	Run:   removeFunc,
}

func init() {
	rootCmd.AddCommand(removeCmd)
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
