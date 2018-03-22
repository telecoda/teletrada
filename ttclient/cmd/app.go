package cmd

import (
	"fmt"
	"log"

	"google.golang.org/grpc/connectivity"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/telecoda/teletrada/proto"
	"google.golang.org/grpc"
)

const (
	defaultAddress = "localhost:13370"
	defaultSymbol  = "btc"
)

var client proto.TeletradaClient
var clientConn *grpc.ClientConn
var address string

var App = grumble.New(&grumble.Config{
	Name:                  "Teletrada",
	Description:           "Crypto trading bot",
	HistoryFile:           "/tmp/teletrada.hist",
	Prompt:                "teletrada » ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.Bool("v", "verbose", false, "enable verbose mode")
		f.String("a", "address", defaultAddress, "address to connect to client")
	},
})

func init() {
	App.SetPrintASCIILogo(func(a *grumble.App) {
		fmt.Println(` _____    _      _                 _       `)
		fmt.Println(`|_   _|  | |    | |               | |      `)
		fmt.Println(`  | | ___| | ___| |_ _ __ __ _  __| | __ _ `)
		fmt.Println(`  | |/ _ \ |/ _ \ __| '__/ _' |/ _' |/ _' |`)
		fmt.Println(`  | |  __/ |  __/ |_| | | (_| | (_| | (_| |`)
		fmt.Println(`  \_/\___|_|\___|\__|_|  \__,_|\__,_|\__,_|`)
		fmt.Println()
	})

	App.OnInit(onInit)
	//flag.StringVar(&address, "address", defaultAddress, "Address to connect to client")

}

func onInit(a *grumble.App, flags grumble.FlagMap) error {
	address = flags.String("address")
	return nil
}

func getClient() proto.TeletradaClient {

	var err error
	if client == nil {
		// Init GRPC client
		fmt.Printf("Connecting to %s\n", address)
		clientConn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to create client connection to: %s - %s", address, err)
		}
		fmt.Printf("Connection successful\n")
		client = proto.NewTeletradaClient(clientConn)
	} else {
		// check current connection state
		state := clientConn.GetState()
		if state != connectivity.Ready {
			printWarningString(fmt.Sprintf("Not connected: current connection state: %d", state))
		}
	}

	return client
}

func Close() {
	// Always close client on exit
	if clientConn != nil {
		if err := clientConn.Close(); err != nil {
			log.Printf("Error closing client conn: %s", err)
		}
	}
}

/*
	TODO: - implement this command hierarchy

	(CRUD) like functionality

	create:
		simulation
		strategy
	list:
		logs
		portfolio
		prices
		simulation(s)
		strategy(ies)
	delete:
		simulation
		strategy
	update:
		simulation
		strategy
	start:
		simulation
	stop:
		simulation
	status:

*/
