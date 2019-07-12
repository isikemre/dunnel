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

// Flag Values
var port int

func handleRootCommand(cmd *cobra.Command, args []string)  {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		network.StartHTTPProxy(port)
		wg.Done()
	}()

	wg.Wait()
}

func initRootCommand() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 1217, "Set the port for dunnel to listen. Default: 1217")
}

func Execute() {
	initRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func AddSubCommands()  {
	rootCmd.AddCommand(VersionCmd)
}