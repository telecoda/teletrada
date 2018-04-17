package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/proto"
)

func TestLogEndpoint(t *testing.T) {

	server, err := setupTestServer()
	assert.NoError(t, err)

	req := &proto.GetLogRequest{}

	rsp, err := server.GetLog(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, rsp)

	// Server startup adds some log entries
	assert.NotEqual(t, 0, len(rsp.Entries))
	logEntriesBefore := len(rsp.Entries)

	DefaultLogger.log("add another message to the log")
	rsp, err = server.GetLog(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.NotEqual(t, logEntriesBefore, len(rsp.Entries))
}

func TestLogPurge(t *testing.T) {

	server, err := setupTestServer()
	assert.NoError(t, err)

	for i := 0; i < MAX_LOG; i++ {
		DefaultLogger.log("add another message to the log")
		// log must get shorter at some point
	}
	req := &proto.GetLogRequest{}
	rsp, err := server.GetLog(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, rsp)

	assert.NotEqual(t, MAX_LOG, len(rsp.Entries))

}
