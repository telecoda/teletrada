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
	header := []string{"sym", "as", "at", "current", "today chg", "today pct", "24h chg", "24h pct", "open", "close", "highest", "lowest", ""}
	writeHeading(tw, header)

	for _, price := range resp.Prices {
		at, err := tspb.Timestamp(price.At)
		if err != nil {
			return err
		}
		writeRow(tw, formatColRow(price.Symbol, price.As, at.Format(DATE_FORMAT), priceField(price.Current), priceField(price.ChangeToday), percentField(price.ChangePctToday), priceField(price.Change24H), percentField(price.ChangePct24H), priceField(price.Opening), priceField(price.Closing), priceField(price.Highest), priceField(price.Lowest), ""))
	}

	tw.Flush()
	fmt.Printf("%s", buf.String())

	return nil
}
