package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dunnel",
	Long:  `All software has versions. This is Dunnel's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dunnel lightweight, blazing-fast and secure Docker tunnel v0.1 -- HEAD")
	},
}

