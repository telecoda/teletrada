package domain

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/telecoda/teletrada/proto"
)

// Rebuild - downloads latest code, recompiles and restarts
func (s *server) Rebuild(ctx context.Context, req *proto.RebuildRequest) (*proto.RebuildResponse, error) {

	resp := &proto.RebuildResponse{}

	log.Printf("Rebuilding ttserver")
	// change dir
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return nil, fmt.Errorf("$GOPATH not found")
	}

	path := filepath.Join(goPath, "src", "github.com", "telecoda", "teletrada", "ttserver")

	if err := os.Chdir(path); err != nil {
		return nil, fmt.Errorf("Failed to change to directory: %s - %s", path, err)
	}

	// pull code
	log.Printf("Pulling latest code...")

	// Untar the new config into the new directory
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// recompile
	log.Printf("Compiling code...")

	// return resp

	// terminate prog..

	resp.Result = "it worked!!!"

	return resp, nil
}
