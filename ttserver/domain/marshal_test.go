package domain

import (
	"testing"
	"time"

	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBalanceAs(t *testing.T) {
	now := time.Now()
	pbNow, _ := tspb.TimestampProto(now)

	tests := []struct {
		balance  *BalanceAs
		expProto *proto.Balance
	}{
		{
			balance: &BalanceAs{
				CoinBalance: exchanges.CoinBalance{
					Symbol:   "symbol",
					Exchange: "exchange",
					Free:     10.0,
					Locked:   20.0,
				},
				Total:        30.0,
				Price:        5.00,
				Value:        150.00,
				At:           now,
				Price24H:     4.00,
				Value24H:     120.00,
				Change24H:    1.00,
				ChangePct24H: 25.00,
			},
			expProto: &proto.Balance{
				Symbol:       "symbol",
				Exchange:     "exchange",
				Free:         10.0,
				Locked:       20.0,
				Total:        30.0,
				Price:        5.00,
				Value:        150.00,
				At:           pbNow,
				Price24H:     4.00,
				Value24H:     120.00,
				Change24H:    1.00,
				ChangePct24H: 25.00,
			},
		},
	}

	for _, test := range tests {
		p, err := test.balance.toProto()
		assert.NoError(t, err)
		assert.Equal(t, test.expProto, p)
	}
}

func TestMarshalLogEntry(t *testing.T) {

	now := time.Now()
	pbNow, _ := tspb.TimestampProto(now)

	tests := []struct {
		entry    *LogEntry
		expProto *proto.LogEntry
	}{
		{
			entry: &LogEntry{
				Timestamp: now,
				Message:   "log message",
			},
			expProto: &proto.LogEntry{
				Time: pbNow,
				Text: "log message",
			},
		},
	}

	for _, test := range tests {
		p, err := test.entry.toProto()
		assert.NoError(t, err)
		assert.Equal(t, test.expProto, p)
	}
}

func TestMarshalPortfolio(t *testing.T) {

	now := time.Now()
	pbNow, _ := tspb.TimestampProto(now)

	tests := []struct {
		portfolio *portfolio
		expProto  *proto.Portfolio
	}{
		{
			portfolio: &portfolio{
				name:   "portfolio-name",
				isLive: true,
				balances: map[SymbolType]*BalanceAs{
					"base-symbol": &BalanceAs{
						CoinBalance: exchanges.CoinBalance{
							Symbol:   "symbol",
							Exchange: "exchange",
							Free:     10.0,
							Locked:   20.0,
						},
						Total:        30.0,
						Price:        5.00,
						Value:        150.00,
						At:           now,
						Price24H:     4.00,
						Value24H:     120.00,
						Change24H:    1.00,
						ChangePct24H: 25.00,
					},
				},
			},
			expProto: &proto.Portfolio{
				Name: "portfolio-name",
				Balances: []*proto.Balance{
					&proto.Balance{
						Symbol:       "symbol",
						Exchange:     "exchange",
						Free:         10.0,
						Locked:       20.0,
						Total:        30.0,
						Price:        5.00,
						Value:        150.00,
						At:           pbNow,
						Price24H:     4.00,
						Value24H:     120.00,
						Change24H:    1.00,
						ChangePct24H: 25.00,
					},
				},
			},
		},
	}

	for _, test := range tests {
		p, err := test.portfolio.toProto()
		assert.NoError(t, err)
		assert.Equal(t, test.expProto, p)
	}
}

func TestMarshalPrice(t *testing.T) {

	now := time.Now()
	pbNow, _ := tspb.TimestampProto(now)

	tests := []struct {
		price    *Price
		expProto *proto.Price
	}{
		{
			price: &Price{
				Base:     "base-symbol",
				Exchange: "exchange",
				As:       "as-symbol",
				At:       now,
				Price:    1234.56,
			},
			expProto: &proto.Price{
				Symbol:   "base-symbol",
				Exchange: "exchange",
				As:       "as-symbol",
				At:       pbNow,
				Current:  1234.56,
			},
		},
	}

	for _, test := range tests {
		p, err := test.price.toProto()
		assert.NoError(t, err)
		assert.Equal(t, test.expProto, p)
	}
}

func TestMarshalSimulation(t *testing.T) {

	now := time.Now()
	pbNow, _ := tspb.TimestampProto(now)

	from := now
	to := now.Add(5 * time.Hour)

	pbFrom, _ := tspb.TimestampProto(from)
	pbTo, _ := tspb.TimestampProto(to)

	tests := []struct {
		simulation *simulation
		expProto   *proto.Simulation
	}{
		{
			simulation: &simulation{
				id:                "sim-id",
				name:              "sim-name",
				isRunning:         true,
				startedTime:       &from,
				stoppedTime:       &to,
				useHistoricalData: true,
				dataFrequency:     time.Duration(1 * time.Minute),
				useRealtimeData:   true,
				simFromTime:       &from,
				simToTime:         &to,
				portfolio: &portfolio{
					name:   "portfolio-name",
					isLive: true,
					balances: map[SymbolType]*BalanceAs{
						"base-symbol": &BalanceAs{
							CoinBalance: exchanges.CoinBalance{
								Symbol:   "symbol",
								Exchange: "exchange",
								Free:     10.0,
								Locked:   20.0,
							},
							Total:        30.0,
							Price:        5.00,
							Value:        150.00,
							At:           now,
							Price24H:     4.00,
							Value24H:     120.00,
							Change24H:    1.00,
							ChangePct24H: 25.00,
						},
					},
				},
			},
			expProto: &proto.Simulation{
				Id:                "sim-id",
				Name:              "sim-name",
				IsRunning:         true,
				StartedTime:       pbFrom,
				StoppedTime:       pbTo,
				UseHistoricalData: true,
				DataFrequency:     60,
				UseRealtimeData:   true,
				FromTime:          pbFrom,
				ToTime:            pbTo,
				Portfolio: &proto.Portfolio{
					Name: "portfolio-name",
					Balances: []*proto.Balance{
						&proto.Balance{
							Symbol:       "symbol",
							Exchange:     "exchange",
							Free:         10.0,
							Locked:       20.0,
							Total:        30.0,
							Price:        5.00,
							Value:        150.00,
							At:           pbNow,
							Price24H:     4.00,
							Value24H:     120.00,
							Change24H:    1.00,
							ChangePct24H: 25.00,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s, err := test.simulation.toProto()
		assert.NoError(t, err)
		assert.Equal(t, test.expProto, s)
	}
}

func TestMarshalStrategy(t *testing.T) {

	tests := []struct {
		strategy *baseStrategy
		expProto *proto.Strategy
	}{
		{
			strategy: &baseStrategy{
				id:          "strategy-id",
				symbol:      SymbolType("base-symbol"),
				as:          SymbolType("as-symbol"),
				coinPercent: 12.34,
				isRunning:   true,
			},
			expProto: &proto.Strategy{
				Id:          "strategy-id",
				Symbol:      "base-symbol",
				As:          "as-symbol",
				CoinPercent: 12.34,
				IsRunning:   true,
				Description: "Base strategy for building other strategies upon",
			},
		},
	}

	for _, test := range tests {
		s, err := strategyToProto(test.strategy)
		assert.NoError(t, err)
		assert.Equal(t, test.expProto, s)
	}
}
