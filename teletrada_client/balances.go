package main

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/abiosoft/ishell"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func getBalances(c *ishell.Context) {
	r, err := client.GetBalances(context.Background(), &proto.BalancesRequest{})
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get balances: %v\n", err)))
		return
	}
	c.Printf("Balances:\n")
	buf := bytes.Buffer{}

	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	header := []string{"sym", "total", "price", "value", "at", ""}
	PrintRow(tw, PaintRowUniformly(GreenText, header))
	PrintRow(tw, PaintRowUniformly(GreenText, AnonymizeRow(header))) // header separator

	for _, balance := range r.Balances {
		PrintRow(tw, FormatRow(balance.Symbol, balance.Total, balance.LatestUSDPrice, -balance.LatestUSDValue, "2006-01-02", ""))
	}

	tw.Flush()
	c.Printf("%s", buf.String())
}
