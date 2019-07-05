package main

import (
	"github.com/isikemre/dunnel/cmd/dunneld/commands"
)

func main() {
	commands.AddSubCommands()
	commands.Execute()
}