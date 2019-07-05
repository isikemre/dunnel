package main

import (
	"github.com/isikemre/dunnel/cmd/dunnel/commands"
)

func main() {
	commands.AddSubCommands()
	commands.Execute()
}