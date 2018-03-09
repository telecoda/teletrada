package cmd

import (
	"fmt"
	"strings"

	"github.com/desertbit/grumble"
)

func init() {
	listCommand := &grumble.Command{
		Name:     "list",
		Aliases:  []string{"ls"},
		Help:     "list <item>",
		LongHelp: "list operations",
	}
	App.AddCommand(listCommand)

	// list portfolio
	listCommand.AddCommand(&grumble.Command{
		Name:      "portfolio",
		Aliases:   []string{"po"},
		Help:      "list portfolio <as>",
		AllowArgs: true,
		Run:       listPortfolio,
	})

	// list prices
	listCommand.AddCommand(&grumble.Command{
		Name:      "prices",
		Aliases:   []string{"pr"},
		Help:      "list prices <base> <as>",
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
}

func listPortfolio(c *grumble.Context) error {
	as := defaultSymbol
	if len(c.Args) >= 1 {
		as = c.Args[0]
	}

	as = strings.ToLower(as)
	if as == "" {
		as = defaultSymbol
	}

	fmt.Printf("Listing portfolio as %q\n", as)
	return nil
}

func listPrices(c *grumble.Context) error {
	base := ""
	as := ""
	if len(c.Args) >= 1 {
		base = c.Args[0]
	}
	if len(c.Args) == 2 {
		as = c.Args[1]
	}

	base = strings.ToLower(base)
	as = strings.ToLower(as)

	allPrices := false
	if base == "" || base == "*all" {
		allPrices = true
	}

	if as == "" {
		as = defaultSymbol
	}

	if allPrices {
		fmt.Printf("Listing all prices as %q\n", as)
	} else {
		fmt.Printf("Listing prices for %q as %q\n", base, as)
	}

	return nil
}

func listSimulations(c *grumble.Context) error {
	fmt.Println("List simulations")
	return nil
}
