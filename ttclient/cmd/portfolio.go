package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/abiosoft/ishell"
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func getPortfolio(c *ishell.Context) {

	req := &proto.GetPortfolioRequest{}

	if len(c.Args) > 0 {
		req.As = c.Args[0]
	}

	r, err := client.GetPortfolio(context.Background(), req)
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get portfolio: %v\n", err)))
		return
	}

	c.Printf("Portfolio balances:\n")

	printBalances(c, r.Balances)
}

func printBalances(c *ishell.Context, balances []*proto.Balance) {
	if len(balances) == 0 {
		return
	}

	buf := bytes.Buffer{}

	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	header := []string{"sym", "as", "total", "price", "price24", "value", "value24", "at", "change24", "changePct", "buystrat", "sellstrat", ""}
	PrintRow(tw, PaintRowUniformly(GreenText, header))
	PrintRow(tw, PaintRowUniformly(GreenText, AnonymizeRow(header))) // header separator

	total := proto.Balance{
		Symbol: "tot",
		As:     balances[0].As,
	}

	for _, balance := range balances {
		at, err := tspb.Timestamp(balance.At)
		if err != nil {
			c.Println(PaintErr(err))
			continue
		}
		buyStrat := ""
		if balance.BuyStrategy != nil {
			buyStrat = balance.BuyStrategy.Id
		}
		sellStrat := ""
		if balance.SellStrategy != nil {
			sellStrat = balance.SellStrategy.Id
		}
		PrintRow(tw, FormatRow(balance.Symbol, balance.As, balance.Total, balance.Price, balance.Price24H, balance.Value, balance.Value24H, at.Format(DATE_FORMAT), balance.Change24H, balance.ChangePct24H, buyStrat, sellStrat, ""))

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
