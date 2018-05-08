package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/ttserver/servertime"
)

func TestPriceValidation(t *testing.T) {

	tests := []struct {
		name        string
		price       Price
		errExpected bool
	}{
		{
			name: "Valid price",
			price: Price{
				Base:  BTC,
				As:    ETH,
				Price: 123.45,
				At:    servertime.Now(),
			},
			errExpected: false,
		},
		{
			name: "Missing time",
			price: Price{
				Base:  BTC,
				As:    ETH,
				Price: 123.45,
			},
			errExpected: true,
		},
		{
			name: "Missing Base",
			price: Price{
				As:    ETH,
				Price: 123.45,
				At:    servertime.Now(),
			},
			errExpected: true,
		},
		{
			name: "Missing As",
			price: Price{
				Base:  BTC,
				Price: 123.45,
				At:    servertime.Now(),
			},
			errExpected: true,
		},
		{
			name: "Missing price",
			price: Price{
				Base: BTC,
				As:   ETH,
				At:   servertime.Now(),
			},
			errExpected: true,
		},
		{
			name: "Negative price",
			price: Price{
				Base:  BTC,
				As:    ETH,
				At:    servertime.Now(),
				Price: -123.45,
			},
			errExpected: true,
		},
		{
			name: "Base==As must == 1.0",
			price: Price{
				Base:  BTC,
				As:    BTC,
				Price: 1.0,
				At:    servertime.Now(),
			},
			errExpected: false,
		},
		// {
		// 	name: "Base==As must != 1.0",
		// 	price: Price{
		// 		Base:  BTC,
		// 		As:    BTC,
		// 		Price: 1.234,
		// 		At:    servertime.Now(),
		// 	},
		// 	errExpected: true,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := test.price.Validate()
			if test.errExpected && err == nil {
				assert.Fail(t, "Error was expected for price")
			}
			if !test.errExpected && err != nil {
				assert.NoError(t, err)
			}
		})
	}

}
