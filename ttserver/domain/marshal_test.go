package domain

import (
	"testing"
	"time"

	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBalance(t *testing.T) {
	now := time.Now()
	pbNow, _ := tspb.TimestampProto(now)

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
				Total: 30.0,
				At:    now,
			},
			expProto: &proto.Balance{
				Symbol:   "symbol",
				Exchange: "exchange",
				Free:     10.0,
				Locked:   20.0,
				Total:    30.0,
				At:       pbNow,
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
