package domain

import (
	"context"
	"log"
	"time"

	"github.com/telecoda/teletrada/proto"
)

type Log struct {
	Timestamp time.Time
	Message   string
}

func (s *server) log(msg string) {
	l := Log{
		Timestamp: time.Now().UTC(),
		Message:   msg,
	}
	s.statusLog = append(s.statusLog, l)
	if s.isVerbose() {
		log.Println(msg)
	}

}

// func (s *server) GetLog() []Log {
// 	return s.statusLog
// }

// GetLog returns server log
func (s *server) GetLog(ctx context.Context, in *proto.LogRequest) (*proto.LogResponse, error) {

	resp := &proto.LogResponse{}

	entries := s.statusLog

	resp.Entries = make([]*proto.LogEntry, len(entries))

	for i, entry := range entries {
		resp.Entries[i] = &proto.LogEntry{
			Time: entry.Timestamp.Format("2006-01-02T03:04:05"),
			Text: entry.Message,
		}
	}

	return resp, nil
}
