package domain

import (
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
)

/* methods here are for convert to/from domain to protobufs */

func (b *Balance) toProto() *proto.Balance {
	return &proto.Balance{
		Symbol:         string(b.Symbol.GetType()),
		Exchange:       b.Exchange,
		Free:           float32(b.Free),
		Locked:         float32(b.Locked),
		Total:          float32(b.Total),
		LatestUSDPrice: float32(b.LatestUSDPrice),
		LatestUSDValue: float32(b.LatestUSDValue),
	}
}

func (l *LogEntry) toProto() (*proto.LogEntry, error) {
	ts, err :=tspb.TimestampProto(l.Timestamp)
	if err !=nil {
		return nil,err
	} 
	return &proto.LogEntry{
		Time: ts,
		Text: l.Message,
	},nil
}
