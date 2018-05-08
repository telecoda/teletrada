package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/proto"
)

func TestStatusEndpoint(t *testing.T) {

	server, err := initMockServer()
	assert.NoError(t, err)

	req := &proto.GetStatusRequest{}

	rsp, err := server.GetStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, rsp)

	assert.Equal(t, int32(3), rsp.TotalSymbols, "Test data should have 3 symbols")
	assert.NotZero(t, rsp.LastUpdate, "Should have updated on init")
	assert.NotZero(t, rsp.ServerStarted, "Should have a start time")
	assert.Equal(t, int32(1), rsp.UpdateCount, "Should have updated prices once at startup")

}
