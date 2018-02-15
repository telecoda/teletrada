package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func getStatus(c *ishell.Context) {
	s, err := client.GetStatus(context.Background(), &proto.StatusRequest{})
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get server status: %v\n", err)))
		return
	}

	c.Printf("Status:\n")
	timestamp, err := tspb.Timestamp(s.ServerStarted)
	if err != nil {
		c.Println(PaintErr(err))
	} else {
		c.Printf("Started: %s\n", timestamp.Format(DATE_FORMAT))
	}
	timestamp, err = tspb.Timestamp(s.LastUpdate)
	if err != nil {
		c.Println(PaintErr(err))
	} else {
		c.Printf("Last update: %s\n", timestamp.Format(DATE_FORMAT))
	}
	c.Printf("Update count: %d\n", s.UpdateCount)
	c.Printf("Total symbols: %d\n", s.TotalSymbols)

}
