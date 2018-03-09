package cmd

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func rebuild(c *ishell.Context) {
	r, err := client.Rebuild(context.Background(), &proto.RebuildRequest{})
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not rebuild code: %v\n", err)))
		return
	}

	c.Printf("Rebuild:\n")
	c.Printf("Result: %s\n", r.Result)
}
