package cmd

import (
	"github.com/desertbit/grumble"
)

func init() {
	startCommand := &grumble.Command{
		Name: "start",
		Help: "start operations",
	}
	App.AddCommand(startCommand)

	// start simulation
	startCommand.AddCommand(&grumble.Command{
		Name:      "simulation",
		Aliases:   []string{"si"},
		Help:      "start simulation",
		Usage:     "start simulation [id]",
		AllowArgs: true,
		Completer: simIdCompleter,
		Run:       startSimulation,
	})
}
