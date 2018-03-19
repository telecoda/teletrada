package main

import (
	"flag"

	"github.com/desertbit/grumble"
	"github.com/telecoda/teletrada/ttclient/cmd"
)

func main() {
	flag.Parse()
	defer cmd.Close()
	grumble.Main(cmd.App)
}
