package domain

import (
	"context"
	"fmt"

	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
)

// GetStatus returns status of server
func (s *server) GetStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {

	s.RLock()
	defer s.RUnlock()

	startTime, err := tspb.TimestampProto(s.startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to convert startTime: %s", err)
	}

	archiveStatus := DefaultArchive.GetStatus()

	lastUpdated, err := tspb.TimestampProto(archiveStatus.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lastUpdate: %s", err)
	}
	resp := &proto.StatusResponse{
		ServerStarted: startTime,
		LastUpdate:    lastUpdated,
		UpdateCount:   int32(archiveStatus.UpdateCount),
		TotalSymbols:  int32(archiveStatus.TotalSymbols),
	}

	return resp, nil
}
