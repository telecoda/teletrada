package main

import (
	"flag"
	"log"

	"github.com/abiosoft/ishell"
	"github.com/telecoda/teletrada/proto"
	"google.golang.org/grpc"
)

const (
	defaultAddress = "localhost:13370"
)

var client proto.TeletradaClient

type params struct {
	address string
}

func (p *params) setup() {
	flag.StringVar(&p.address, "address", defaultAddress, "Address to connect to client")
}

func main() {
	p := &params{}
	p.setup()
	flag.Parse()

	// Init GRPC client
	conn, err := grpc.Dial(p.address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = proto.NewTeletradaClient(conn)

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println(Paint(YellowText, "Welcome to the TeleTrada client"))

	registerCmds(shell)

	// Read and write history to $HOME/.teletrada_history
	shell.SetHomeHistoryPath(".teletrada_history")

	// run shell
	shell.Run()
}

func registerCmds(shell *ishell.Shell) {
	// register a functions commands.
	shell.AddCmd(&ishell.Cmd{
		Name: "balances",
		Help: "show current balances eg. balances BTC",
		Func: getBalances,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "log",
		Help: "show server logs",
		Func: getLog,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "rebuild",
		Help: "rebuild with latest code and restart",
		Func: rebuild,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "status",
		Help: "show server status",
		Func: getStatus,
	})

}
