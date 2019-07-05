package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dunnel",
	Short: "Dunnel is a lightweight, fast and secure tunnel between a Docker client and server",
	Long: `It's a simple spell, but quite unbreakable.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func AddSubCommands()  {
	rootCmd.AddCommand(VersionCmd)
}