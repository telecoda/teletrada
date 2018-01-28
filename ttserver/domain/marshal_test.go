package domain

import (
	"testing"
	"time"

	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBalance(t *testing.T) {
	tests := []struct {
		balance  *Balance
		expProto *proto.Balance
	}{
		{
			balance: &Balance{
				ExchangeBalance: exchanges.ExchangeBalance{
					Symbol:   "symbol",
					Exchange: "exchange",
					Free:     10.0,
					Locked:   20.0,
				},
				Symbol: &symbol{
					SymbolType: SymbolType("symbol"),
				},
				Total:          30.0,
				Value:          1234.56,
				LatestUSDPrice: 1234.56,
				LatestUSDValue: 30 * 1234.56,
			},
			expProto: &proto.Balance{
				Symbol:         "symbol",
				Exchange:       "exchange",
				Free:           10.0,
				Locked:         20.0,
				Total:          30.0,
				LatestUSDPrice: 1234.56,
				LatestUSDValue: 30 * 1234.56,
			},
		},
	}

	for _, test := range tests {
		p := test.balance.toProto()
		assert.Equal(t, test.expProto, p)
	}
}

func TestMarshalLogEntry(t *testing.T) {

	now := time.Now()

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
				Time: &google_protobuf.Timestamp{Nanos: int32(now.Nanosecond())},
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
