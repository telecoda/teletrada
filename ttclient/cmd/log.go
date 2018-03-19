package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/desertbit/grumble"
	tspb "github.com/golang/protobuf/ptypes"
	"github.com/telecoda/teletrada/proto"
	"golang.org/x/net/context"
)

const DATE_FORMAT = "2006-01-02 15:04:05"

func listLogs(c *grumble.Context) error {
	r, err := getClient().GetLog(context.Background(), &proto.GetLogRequest{})
	if err != nil {
		return fmt.Errorf("could not get server log: %v\n", err)
	}

	buf := bytes.Buffer{}
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Headers
	header := []string{"timestamp", "log", ""}
	writeHeading(tw, header)

	for _, entry := range r.Entries {
		timestamp, err := tspb.Timestamp(entry.Time)
		if err != nil {
			return err
		}
		writeRow(tw, formatColRow(timestamp.Format(DATE_FORMAT), entry.Text, ""))
	}
	tw.Flush()
	fmt.Printf("%s", buf.String())

	return nil
}
