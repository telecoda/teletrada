package cmd

import (
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func init() {
	App.AddCommand(&grumble.Command{
		Name: "rebuild",
		Help: "rebuild server and restart",
		Run:  rebuild,
	})
}

func rebuild(c *grumble.Context) error {

	r, err := getClient().Rebuild(context.Background(), &proto.RebuildRequest{})
	if err != nil {
		return (fmt.Errorf("could not rebuild code: %v\n", err))
	}

	printHeading("Rebuild")
	fmt.Printf("Result: %s\n", r.Result)

	return nil
}
