package cmd

import (
	"sort"
	"strings"

	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

var symbolTypes map[string][]string
var symbols []string

func initSymbolTypes() {
	resp, err := getClient().GetSymbolTypes(context.Background(), &proto.GetSymbolTypesRequest{})
	if err != nil {
		return
	}

	// convert symbols to list for completer
	symbols = make([]string, len(resp.SymbolTypes))
	symbolTypes = make(map[string][]string, len(symbols))
	i := 0
	for _, symbolType := range resp.SymbolTypes {
		symbolType.Base = strings.ToLower(symbolType.Base)

		for i, as := range symbolType.As {
			symbolType.As[i] = strings.ToLower(as)
		}
		symbolTypes[symbolType.Base] = symbolType.As
		symbols[i] = symbolType.Base
		i++
	}
	sort.Strings(symbols)
}

func symbolCompleter(prefix string, args []string) []string {

	if len(symbolTypes) == 0 {
		initSymbolTypes()
	}

	if len(args) == 0 {
		if prefix == "" {
			// return top level symbols
			return symbols
		} else {
			// only matching suffix
			suffixes := make([]string, 0)
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, prefix) {
					suffixes = append(suffixes, symbol)
				}
			}
			return suffixes
		}
	}
	if len(args) == 1 && prefix == "" {
		// check for space between next arg
		// symbol := args[0]
		as := symbolTypes[args[0]]
		sort.Strings(as)
		return as
	}
	return []string{}
}
