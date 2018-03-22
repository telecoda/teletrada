package main

import (
	"github.com/desertbit/grumble"
	"github.com/telecoda/teletrada/ttclient/cmd"
)

func main() {
	defer cmd.Close()
	grumble.Main(cmd.App)
}
