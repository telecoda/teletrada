package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func getSimulations(c *ishell.Context) {

	req := &proto.GetSimulationsRequest{}

	if len(c.Args) > 0 {
		req.Id = c.Args[0]
	}

	r, err := client.GetSimulations(context.Background(), req)
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get simulations: %v\n", err)))
		return
	}

	c.Printf("Simulations:\n")
	c.Printf("============\n")

	for _, simulation := range r.Simulations {
		if simulation.Portfolio != nil {
			// Print simulation header details
			c.Printf("Name: %s\n", simulation.Portfolio.Name)
			// Print simulation balances

			printBalances(c, simulation.Portfolio.Balances)
		}
	}
}
