package domain

import (
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
)

/* methods here are for convert to/from domain to protobufs */

func (b *BalanceAs) toProto() (*proto.Balance, error) {
	pb := &proto.Balance{
		Symbol:       b.CoinBalance.Symbol,
		Exchange:     b.Exchange,
		Free:         float32(b.Free),
		Locked:       float32(b.Locked),
		Total:        float32(b.Total),
		As:           string(b.As),
		Price:        float32(b.Price),
		Value:        float32(b.Value),
		Price24H:     float32(b.Price24H),
		Value24H:     float32(b.Value24H),
		Change24H:    float32(b.Change24H),
		ChangePct24H: float32(b.ChangePct24H),
	}

	ts, err := tspb.TimestampProto(b.At)
	if err != nil {
		return nil, err
	} else {
		pb.At = ts
	}

	if b.BuyStrategy != nil {
		if pb.BuyStrategy, err = b.BuyStrategy.toProto(); err != nil {
			return nil, err
		}
	}

	if b.SellStrategy != nil {
		if pb.SellStrategy, err = b.SellStrategy.toProto(); err != nil {
			return nil, err
		}
	}

	return pb, nil
}

func (l *LogEntry) toProto() (*proto.LogEntry, error) {
	ts, err := tspb.TimestampProto(l.Timestamp)
	if err != nil {
		return nil, err
	}
	return &proto.LogEntry{
		Time: ts,
		Text: l.Message,
	}, nil
}

func (p *portfolio) toProto() (*proto.Portfolio, error) {
	pp := &proto.Portfolio{
		Name:     p.name,
		Balances: make([]*proto.Balance, len(p.balances)),
	}

	i := 0
	var err error
	for _, balance := range p.balances {
		if pp.Balances[i], err = balance.toProto(); err != nil {
			return nil, err
		}
		i++
	}

	return pp, nil
}

func (p *Price) toProto() (*proto.Price, error) {
	pp := &proto.Price{
		Symbol:   string(p.Base),
		Exchange: p.Exchange,
		As:       string(p.As),
		Current:  float32(p.Price),
	}

	ts, err := tspb.TimestampProto(p.At)
	if err != nil {
		return nil, err
	} else {
		pp.At = ts
	}

	return pp, nil
}

func (s *simulation) toProto() (*proto.Simulation, error) {
	ps := &proto.Simulation{
		UseHistoricalData: s.useHistoricData,
		DataFrequency:     int32(s.dataFrequency.Seconds()),
		UseRealtimeData:   s.useRealtimeData,
	}

	if s.simFromTime != nil {
		fromTime, err := tspb.TimestampProto(*s.simFromTime)
		if err != nil {
			return nil, err
		}
		ps.FromTime = fromTime
	}

	if s.simToTime != nil {
		toTime, err := tspb.TimestampProto(*s.simToTime)
		if err != nil {
			return nil, err
		}
		ps.ToTime = toTime
	}

	if protoPort, err := s.portfolio.toProto(); err != nil {
		return nil, err
	} else {
		ps.Portfolio = protoPort
	}

	return ps, nil
}

func (b *baseStrategy) toProto() (*proto.Strategy, error) {
	ps := &proto.Strategy{
		Id:          b.ID(),
		Description: b.Description(),
		CoinPercent: float32(b.CoinPercent()),
		Symbol:      string(b.Symbol()),
		As:          string(b.As()),
		IsRunning:   b.IsRunning(),
	}
	return ps, nil
}
