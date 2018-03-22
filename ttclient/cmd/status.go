package cmd

import (
	"fmt"
	"time"

	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"

	"github.com/desertbit/grumble"
)

func init() {
	App.AddCommand(&grumble.Command{
		Name:      "status",
		Help:      "show server status",
		Aliases:   []string{"st"},
		Usage:     "status",
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Run: getStatus,
	})
}

func getStatus(c *grumble.Context) error {
	printHeading("Status")
	fmt.Print(formatAttrString("Address", address+"\n"))
	s, err := getClient().GetStatus(context.Background(), &proto.GetStatusRequest{})
	if err != nil {
		return err
	}

	timestamp, err := tspb.Timestamp(s.ServerStarted)
	if err != nil {
		return err
	} else {
		fmt.Print(formatAttrString("Started", timestamp.Format(DATE_FORMAT)) + "\n")
	}
	timestamp, err = tspb.Timestamp(s.LastUpdate)
	if err != nil {
		return err
	} else {
		fmt.Print(formatAttrString("Last updated", timestamp.Format(DATE_FORMAT)) + "\n")
	}
	fmt.Print(formatAttrInt("Update count", int(s.UpdateCount)) + "\n")
	fmt.Print(formatAttrInt("Total symbols", int(s.TotalSymbols)) + "\n")

	return nil

}
