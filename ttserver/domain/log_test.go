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
	assert.Equal(t, 3, len(rsp.Entries))

	server.log("add another message to the log")

	rsp, err = server.GetLog(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, 4, len(rsp.Entries))
}
