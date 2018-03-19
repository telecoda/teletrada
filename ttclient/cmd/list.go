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

func init() {
	listCommand := &grumble.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Help:    "list operations",
	}
	App.AddCommand(listCommand)

	// list logs
	listCommand.AddCommand(&grumble.Command{
		Name: "logs",
		Help: "list logs",
		Run:  listLogs,
	})
	// list portfolio
	listCommand.AddCommand(&grumble.Command{
		Name:      "portfolio",
		Aliases:   []string{"po"},
		Help:      "list portfolio",
		Usage:     "list portfolio [as]",
		AllowArgs: true,
		Run:       listPortfolio,
	})

	// list prices
	listCommand.AddCommand(&grumble.Command{
		Name:      "prices",
		Aliases:   []string{"pr"},
		Help:      "list prices",
		Usage:     "list prices [base] [as]",
		AllowArgs: true,
		Completer: symbolCompleter,
		Run:       listPrices,
	})

	// list simulations
	listCommand.AddCommand(&grumble.Command{
		Name:    "simulations",
		Aliases: []string{"si"},
		Help:    "list simulations",
		Run:     listSimulations,
	})
}

func listPrices(c *grumble.Context) error {
	base := ""
	as := ""
	if len(c.Args) >= 1 {
		base = c.Args[0]
	}
	if len(c.Args) == 2 {
		as = c.Args[1]
	}

	base = strings.ToLower(base)
	as = strings.ToLower(as)

	allPrices := false
	if base == "" || base == "*all" {
		allPrices = true
	}

	if as == "" {
		as = defaultSymbol
	}

	if allPrices {
		fmt.Printf("Listing all prices as %q\n", as)
	} else {
		fmt.Printf("Listing prices for %q as %q\n", base, as)
	}

	resp, err := getClient().GetPrices(context.Background(), &proto.GetPricesRequest{Base: base, As: as})
	if err != nil {
		return err
	}

	// print prices
	printHeading("Prices")

	buf := bytes.Buffer{}

	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	header := []string{"sym", "as", "price", "price24", "at", "change24", "changePct", ""}
	writeHeading(tw, header)

	for _, price := range resp.Prices {
		at, err := tspb.Timestamp(price.At)
		if err != nil {
			return err
		}
		writeRow(tw, formatColRow(price.Symbol, price.As, priceField(price.Price), priceField(price.Price24H), at.Format(DATE_FORMAT), priceField(price.Change24H), percentField(price.ChangePct24H), ""))
	}

	tw.Flush()
	fmt.Printf("%s", buf.String())

	return nil
}
