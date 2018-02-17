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
