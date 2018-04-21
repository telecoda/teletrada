package cmd

import (
	"fmt"
	"time"

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
			printHeading(fmt.Sprintf("Id: %s", simulation.Id))
			fmt.Print(formatAttrString("Name",
				simulation.Name+"\n"))
			fmt.Print(formatAttrString("Data freq",
				(time.Duration(int64(simulation.DataFrequency))*time.Second).String()+"\n"))
			fmt.Print(formatAttrString("Use historical",
				fmt.Sprintf("%t", simulation.UseHistoricalData)+"\n"))
			if simulation.UseHistoricalData {
				// show times as well
				fmt.Print(formatAttrString("Hist From", formatProtoTimestamp(simulation.FromTime)+"\n"))
				fmt.Print(formatAttrString("Hist To  ", formatProtoTimestamp(simulation.ToTime)+"\n"))
			}
			fmt.Print(formatAttrString("Use realtime",
				fmt.Sprintf("%t", simulation.UseRealtimeData)+"\n"))
			fmt.Print(formatAttrString("IsRunning",
				fmt.Sprintf("%t", simulation.IsRunning)+"\n"))
			fmt.Print(formatAttrString("Started", formatProtoTimestamp(simulation.StartedTime)+"\n"))
			fmt.Print(formatAttrString("Stopped", formatProtoTimestamp(simulation.StoppedTime)+"\n"))

			printBalances(simulation.Portfolio.Balances)
		}
	}

	return nil
}

func createSimulation(c *grumble.Context) error {

	printHeading("Create simulation")
	req := &proto.CreateSimulationRequest{}

	r, err := getClient().CreateSimulation(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create simulation: %v\n", err)
	}

	fmt.Printf("Resp: %#v\n", r)
	return nil
}

func startSimulation(c *grumble.Context) error {

	printHeading("Start simulation")

	if len(c.Args) != 1 {
		return fmt.Errorf("you must provide a simulation id")
	}
	req := &proto.StartSimulationRequest{
		Id: c.Args[0],
	}

	_, err := getClient().StartSimulation(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to start simulation: %v\n", err)
	}

	fmt.Printf("Simulation started\n")
	return nil
}

func stopSimulation(c *grumble.Context) error {

	printHeading("Stop simulation")
	if len(c.Args) != 1 {
		return fmt.Errorf("you must provide a simulation id")
	}
	req := &proto.StopSimulationRequest{
		Id: c.Args[0],
	}

	_, err := getClient().StopSimulation(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to stop simulation: %v\n", err)
	}

	fmt.Printf("Simulation stopped\n")
	return nil
}
