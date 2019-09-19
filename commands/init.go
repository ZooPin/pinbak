package commands

import (
	"github.com/ZooPin/pinbak/helper"
	"github.com/ZooPin/pinbak/manager"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init pinbak.",
	Run:   initFunc,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initFunc(cmd *cobra.Command, args []string) {
	_, err := helper.GetConfig()
	if err == nil {
		return
	}

	err = os.Mkdir(helper.PinbakPath(), 0766)
	if err != nil {
		log.Fatal("Init error: ", err)
	}

	var config manager.Config
	config.Name = "Pinbak"
	config.Email = "no-email@pinbak"
	config.SetPath(helper.PinbakPath())
	err = config.Save()
	if err != nil {
		log.Fatal("Init error: ", err)
	}
}
