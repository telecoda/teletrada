package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/abiosoft/ishell"
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

const DATE_FORMAT = "2006-01-02 15:04:05"

func getLog(c *ishell.Context) {
	r, err := client.GetLog(context.Background(), &proto.GetLogRequest{})
	if err != nil {
		c.Print(PaintErr(fmt.Errorf("could not get server log: %v\n", err)))
		return
	}

	buf := bytes.Buffer{}
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Headers
	header := []string{"timestamp", "log", ""}
	PrintRow(tw, PaintRowUniformly(GreenText, header))
	PrintRow(tw, PaintRowUniformly(GreenText, AnonymizeRow(header))) // header separator

	for _, entry := range r.Entries {
		timestamp, err := tspb.Timestamp(entry.Time)
		if err != nil {
			c.Println(PaintErr(err))
			continue
		}
		PrintRow(tw, FormatRow(timestamp.Format(DATE_FORMAT), entry.Text, ""))
	}
	tw.Flush()
	c.Printf("%s", buf.String())
}
