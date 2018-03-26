package domain

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		output, _ := cmd.Output()
		return nil, fmt.Errorf("Failed fetch latest code %s - %s", err, string(output))
	}
	output, _ := cmd.Output()
	log.Printf("Fetched latest code... %s", string(output))

	// recompile
	log.Printf("Compiling code...")
	cmd = exec.Command("go", "install", "github.com/telecoda/teletrada/ttserver", "-a")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		output, _ := cmd.Output()
		return nil, fmt.Errorf("Failed compile latest code %s - %s", err, string(output))
	}

	// terminate prog..
	go func() {
		// sleep then die
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	resp.Result = "Rebuild successful"

	return resp, nil
}
