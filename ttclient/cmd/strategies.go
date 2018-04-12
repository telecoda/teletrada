package cmd

import (
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func listStrategies(c *grumble.Context) error {

	printHeading("List strategies")

	// for now just fetch the strategies on the simulation
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
			printHeading(fmt.Sprintf("Name: %s PortName: %s", simulation.Name, simulation.Portfolio.Name))

			for _, balance := range simulation.Portfolio.Balances {
				if balance.BuyStrategy != nil {
					fmt.Printf("Sim: %s BuyStrat: %s\n", simulation.Name, balance.BuyStrategy.Description)
				}
				if balance.BuyStrategy != nil {
					fmt.Printf("Sim: %s SellStrat: %s\n", simulation.Name, balance.SellStrategy.Description)
				}
			}
		}
	}

	// TODO: add code to fetch strategies from the live portfolio too
	return nil
}
