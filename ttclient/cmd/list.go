package cmd

import (
	"github.com/desertbit/grumble"
)

func init() {
	listCommand := &grumble.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Help:    "list operations",
	}
	App.AddCommand(listCommand)

	// list logs
	listCommand.AddCommand(&grumble.Command{
		Name:    "logs",
		Aliases: []string{"lo"},
		Help:    "list logs",
		Run:     listLogs,
	})

	// list status
	listCommand.AddCommand(&grumble.Command{
		Name:    "status",
		Aliases: []string{"sts"},
		Help:    "show server status",
		Run:     getStatus,
	})

	// list portfolio
	listCommand.AddCommand(&grumble.Command{
		Name:      "portfolio",
		Aliases:   []string{"po"},
		Help:      "list portfolio",
		Usage:     "list portfolio [as]",
		AllowArgs: true,
		Run:       listPortfolio,
	})

	// list prices
	listCommand.AddCommand(&grumble.Command{
		Name:      "prices",
		Aliases:   []string{"pr"},
		Help:      "list prices",
		Usage:     "list prices [base] [as]",
		AllowArgs: true,
		Completer: symbolCompleter,
		Run:       listPrices,
	})

	// list simulations
	listCommand.AddCommand(&grumble.Command{
		Name:    "simulations",
		Aliases: []string{"si"},
		Help:    "list simulations",
		Run:     listSimulations,
	})

	// list strategies
	listCommand.AddCommand(&grumble.Command{
		Name:    "strategies",
		Aliases: []string{"st"},
		Help:    "list strategies",
		Run:     listStrategies,
	})

}
