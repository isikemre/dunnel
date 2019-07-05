package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dunneld",
	Long:  `All software has versions. This is Dunneld's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dunneld is the daemon which will be installed on the remote host for dunnel with docker installed v0.1 -- HEAD")
	},
}

