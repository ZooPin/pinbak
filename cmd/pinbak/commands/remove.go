package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pinbak/cmd/pinbak/helper"
)

var removeCmd = &cobra.Command{
	Use:   "remove [repository] [id] [id...]",
	Short: "Remove an item from the backup.",
	Args:  cobra.MinimumNArgs(2),
	Run:   removeFunc,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeFunc(cmd *cobra.Command, args []string) {
	mover, err := helper.GetMover()
	if err != nil {
		log.Fatal("Error remove: ", err)
	}

	for i := 1; i < len(args); i++ {
		err = mover.Remove(args[0], args[i])
		if err != nil {
			log.Print("Remove error in ", args[0], " file ", args[i], ": ", err)
		}
	}
	err = mover.Git.CommitAndPush(args[0])
	if err != nil {
		log.Fatalln("Remove error: ", err)
	}

	fmt.Println("Done.")
}
