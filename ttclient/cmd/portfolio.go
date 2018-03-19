package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/desertbit/grumble"
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

func listPortfolio(c *grumble.Context) error {
	as := defaultSymbol
	if len(c.Args) >= 1 {
		as = c.Args[0]
	}

	as = strings.ToLower(as)
	if as == "" {
		as = defaultSymbol
	}

	printHeading(fmt.Sprintf("Listing portfolio as %q", as))
	req := &proto.GetPortfolioRequest{
		As: as,
	}

	r, err := getClient().GetPortfolio(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not get portfolio: %v\n", err)
	}

	if len(r.Balances) == 0 {
		return nil
	}

	return printBalances(r.Balances)
}

func printBalances(balances []*proto.Balance) error {
	buf := bytes.Buffer{}

	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	header := []string{"sym", "as", "total", "price", "price24", "value", "value24", "at", "change24", "changePct", "buystrat", "sellstrat", ""}
	writeHeading(tw, header)

	total := proto.Balance{
		Symbol: "tot",
		As:     balances[0].As,
	}

	for _, balance := range balances {
		at, err := tspb.Timestamp(balance.At)
		if err != nil {
			return err
		}
		buyStrat := ""
		if balance.BuyStrategy != nil {
			buyStrat = balance.BuyStrategy.Id
		}
		sellStrat := ""
		if balance.SellStrategy != nil {
			sellStrat = balance.SellStrategy.Id
		}
		writeRow(tw, formatColRow(balance.Symbol, balance.As, priceField(balance.Total), priceField(balance.Price), priceField(balance.Price24H), priceField(balance.Value), priceField(balance.Value24H), at.Format(DATE_FORMAT), priceField(balance.Change24H), percentField(balance.ChangePct24H), buyStrat, sellStrat, ""))

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
	writeRow(tw, formatColRow(total.Symbol, total.As, "", "", "", priceField(total.Value), priceField(total.Value24H), "", priceField(total.Change24H), priceField(total.ChangePct24H), ""))

	tw.Flush()
	fmt.Printf("%s", buf.String())

	return nil
}
