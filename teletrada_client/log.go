package main

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/abiosoft/ishell"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func getLog(c *ishell.Context) {
	r, err := client.GetLog(context.Background(), &proto.LogRequest{})
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get server log: %v\n", err)))
		return
	}

	buf := bytes.Buffer{}
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Headers
	header := []string{"timestamp", "log", ""}
	PrintRow(tw, PaintRowUniformly(GreenText, header))
	PrintRow(tw, PaintRowUniformly(GreenText, AnonymizeRow(header))) // header separator

	for _, entry := range r.Entries {
		PrintRow(tw, FormatRow(entry.Time, entry.Text, ""))
	}
	tw.Flush()
	c.Printf("%s", buf.String())
}
