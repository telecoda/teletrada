package cmd

import (
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func listSimulations(c *grumble.Context) error {

	printHeading("List simulations")
	req := &proto.GetSimulationsRequest{}

	if len(c.Args) > 0 {
		req.Id = c.Args[0]
	}

	r, err := getClient().GetSimulations(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not get simulations: %v\n", err)
	}

	for _, simulation := range r.Simulations {
		if simulation.Portfolio != nil {
			// Print simulation header details
			printHeading(fmt.Sprintf("Name: %s", simulation.Portfolio.Name))

			printBalances(simulation.Portfolio.Balances)
		}
	}

	return nil
}
