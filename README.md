# teletrada
Crypto trading bot, consisting a gRPC client and server.

Very WIP at this stage. Very little will work at present.

# ttserver

This is the server the communicates with the exchange(s) and accumulates pricing information to simulating trading stategies.

# ttclient

This is the client that communicates with the server.  The client has an interactive shell to allow you to query and control the server component.

## Installation

    go get github.com/telecoda/teletrada/ttserver

## Running tests

   go test ./...


## Dependencies

- Install influxdb

    brew install influxdb

    brew services start influxdb

- Install Chronograph

    brew install yarn --without-node

    go get github.com/influxdata/chronograf
    cd $GOPATH/src/github.com/influxdata/chronograf
    make
    