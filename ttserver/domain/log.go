package domain

import (
	"context"
	"log"
	"time"

	"github.com/telecoda/teletrada/proto"
)

var MAX_LOG = 1000

var DefaultLogger = NewLogger(false)

type Logger interface {
	log(message string)
	GetEntries() []LogEntry
}

type logger struct {
	isVerbose bool
	// logging
	statusLog []LogEntry
}

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

func NewLogger(isVerbose bool) Logger {
	logger := &logger{
		isVerbose: isVerbose,
	}

	return logger
}

func (l *logger) log(msg string) {
	entry := LogEntry{
		Timestamp: ServerTime(),
		Message:   msg,
	}
	l.statusLog = append(l.statusLog, entry)
	if l.isVerbose {
		log.Println(msg)
	}

	if len(l.statusLog) == MAX_LOG {
		// purge the log (save last half)
		l.statusLog = l.statusLog[MAX_LOG/2 : MAX_LOG]
	}

}

func (l *logger) GetEntries() []LogEntry {
	return l.statusLog
}

// GetLog returns server log
func (s *server) GetLog(ctx context.Context, in *proto.GetLogRequest) (*proto.GetLogResponse, error) {

	resp := &proto.GetLogResponse{}

	entries := DefaultLogger.GetEntries()

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
