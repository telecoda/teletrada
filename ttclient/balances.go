package main

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/abiosoft/ishell"
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func getBalances(c *ishell.Context) {

	req := &proto.BalancesRequest{}

	if len(c.Args) > 0 {
		req.As = c.Args[0]
	}

	r, err := client.GetBalances(context.Background(), req)
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get balances: %v\n", err)))
		return
	}

	c.Printf("Balances:\n")
	buf := bytes.Buffer{}

	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	header := []string{"sym", "as", "total", "price", "value", "at", ""}
	PrintRow(tw, PaintRowUniformly(GreenText, header))
	PrintRow(tw, PaintRowUniformly(GreenText, AnonymizeRow(header))) // header separator

	for _, balance := range r.Balances {
		at, err := tspb.Timestamp(balance.At)
		if err != nil {
			c.Println(PaintErr(err))
			continue
		}
		PrintRow(tw, FormatRow(balance.Symbol, balance.As, balance.Total, balance.AsPrice, balance.AsValue, at.Format(DATE_FORMAT), ""))
	}

	tw.Flush()
	c.Printf("%s", buf.String())
}
