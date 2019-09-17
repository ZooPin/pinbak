package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "pinbak",
	Short: "pinbak is a small backup manager",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func homeDir() string {
	s, _ := os.UserHomeDir()
	return s
}

func pinBakPath() string {
	return fmt.Sprint(homeDir(), "/.pinbak")
}

func configPath() string {
	return fmt.Sprint(pinBakPath(), "/config")
}
