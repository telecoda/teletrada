package domain

import (
	"context"
	"log"
	"time"

	"github.com/telecoda/teletrada/proto"
)

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

func (s *server) log(msg string) {
	l := LogEntry{
		Timestamp: time.Now().UTC(),
		Message:   msg,
	}
	s.statusLog = append(s.statusLog, l)
	if s.isVerbose() {
		log.Println(msg)
	}

}

// GetLog returns server log
func (s *server) GetLog(ctx context.Context, in *proto.LogRequest) (*proto.LogResponse, error) {

	resp := &proto.LogResponse{}

	entries := s.statusLog

	resp.Entries = make([]*proto.LogEntry, len(entries))

	var err error
	for i, entry := range entries {
		resp.Entries[i], err = entry.toProto()
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
