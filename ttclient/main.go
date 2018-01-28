package main

import (
	"log"

	"github.com/abiosoft/ishell"
	"github.com/telecoda/teletrada/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var client proto.TeletradaClient

func main() {

	// Init GRPC client
	conn, err := grpc.Dial(address, grpc.WithInsecure())
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
	// register a function for "balances" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "balances",
		Help: "show current balances",
		Func: getBalances,
	})

	// register a function for "log" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "log",
		Help: "show server logs",
		Func: getLog,
	})

}
