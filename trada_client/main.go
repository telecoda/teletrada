package main

import (
	"log"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var client proto.TraderClient

func main() {

	// Init GRPC client
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = proto.NewTraderClient(conn)

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("Welcome to the TeleTrada client")

	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "getbalances",
		Help: "get current balances",
		Func: getBalances,
	})
	// Read and write history to $HOME/.teletrada_history
	shell.SetHomeHistoryPath(".teletrada_history")

	// run shell
	shell.Run()
}

func getBalances(c *ishell.Context) {
	r, err := client.GetBalances(context.Background(), &proto.BalancesRequest{})
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		c.Printf(red("could not get balances: %v\n", err))
		return
	}
	c.Printf("Balances: %s\n", r.Message)
}
