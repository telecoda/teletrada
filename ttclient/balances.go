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
	header := []string{"sym", "as", "total", "price", "price24", "value", "value24", "at", "change24", "changePct", ""}
	PrintRow(tw, PaintRowUniformly(GreenText, header))
	PrintRow(tw, PaintRowUniformly(GreenText, AnonymizeRow(header))) // header separator

	total := proto.Balance{
		Symbol: "tot",
		As:     req.As,
	}

	for _, balance := range r.Balances {
		at, err := tspb.Timestamp(balance.At)
		if err != nil {
			c.Println(PaintErr(err))
			continue
		}
		PrintRow(tw, FormatRow(balance.Symbol, balance.As, balance.Total, balance.Price, balance.Price24H, balance.Value, balance.Value24H, at.Format(DATE_FORMAT), balance.Change24H, balance.ChangePct24H, ""))

		// add to total
		total.Exchange = balance.Exchange
		total.Value += balance.Value
		total.Value24H += balance.Value24H
	}

	// calc change
	total.Change24H = total.Value - total.Value24H
	if total.Value != 0 {
		total.ChangePct24H = total.Change24H / total.Value
	}

	// Print total
	PrintRow(tw, FormatRow(total.Symbol, total.As, "", "", "", total.Value, total.Value24H, "", total.Change24H, total.ChangePct24H, ""))

	tw.Flush()
	c.Printf("%s", buf.String())
}
