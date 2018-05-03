package cmd

import (
	"github.com/desertbit/grumble"
)

func init() {
	createCommand := &grumble.Command{
		Name:    "create",
		Aliases: []string{"cr"},
		Help:    "create operations",
	}
	App.AddCommand(createCommand)

	// create simulation
	createCommand.AddCommand(&grumble.Command{
		Name:      "simulation",
		Aliases:   []string{"si"},
		Help:      "create simulation",
		Usage:     "create simulation [id] [name]",
		AllowArgs: true,
		Run:       createSimulation,
	})
}
