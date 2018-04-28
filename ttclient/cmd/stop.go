package cmd

import (
	"github.com/desertbit/grumble"
)

func init() {
	stopCommand := &grumble.Command{
		Name: "stop",
		Help: "stop operations",
	}
	App.AddCommand(stopCommand)

	// stop simulation
	stopCommand.AddCommand(&grumble.Command{
		Name:      "simulation",
		Aliases:   []string{"si"},
		Help:      "stop simulation",
		Usage:     "stop simulation [id]",
		AllowArgs: true,
		Completer: simStopCompleter,
		Run:       stopSimulation,
	})
}
