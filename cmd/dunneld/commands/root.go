package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/isikemre/dunnel/pkg/network"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dunneld",
	Short: "Dunneld is the daemon which will be installed on the remote host for dunnel with docker installed",
	Long: `It's a simple spell, but quite unbreakable.`,
	Run: handleRootCommand,
}

func handleRootCommand(cmd *cobra.Command, args []string)  {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		network.StartHTTPProxy()
		wg.Done()
	}()

	go func() {
		network.StartTCPProxy()
		wg.Done()
	}()

	wg.Wait()
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