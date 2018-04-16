package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	tspb "github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
)

type Color string

// Color codes interpretted by the terminal
// NOTE: all codes must be of the same length or they will throw off the field alignment of tabwriter
const (
	Reset                   Color = "\x1b[0000m"
	Bright                        = "\x1b[0001m"
	BlackText                     = "\x1b[0030m"
	RedText                       = "\x1b[0031m"
	GreenText                     = "\x1b[0032m"
	YellowText                    = "\x1b[0033m"
	BlueText                      = "\x1b[0034m"
	MagentaText                   = "\x1b[0035m"
	CyanText                      = "\x1b[0036m"
	WhiteText                     = "\x1b[0037m"
	DefaultText                   = "\x1b[0039m"
	BrightRedText                 = "\x1b[1;31m"
	BrightGreenText               = "\x1b[1;32m"
	BrightYellowText              = "\x1b[1;33m"
	BrightBlueText                = "\x1b[1;34m"
	BrightMagentaText             = "\x1b[1;35m"
	BrightCyanText                = "\x1b[1;36m"
	BrightWhiteText               = "\x1b[1;37m"
	BlackBackground               = "\x1b[0040m"
	RedBackground                 = "\x1b[0041m"
	GreenBackground               = "\x1b[0042m"
	YellowBackground              = "\x1b[0043m"
	BlueBackground                = "\x1b[0044m"
	MagentaBackground             = "\x1b[0045m"
	CyanBackground                = "\x1b[0046m"
	WhiteBackground               = "\x1b[0047m"
	BrightBlackBackground         = "\x1b[0100m"
	BrightRedBackground           = "\x1b[0101m"
	BrightGreenBackground         = "\x1b[0102m"
	BrightYellowBackground        = "\x1b[0103m"
	BrightBlueBackground          = "\x1b[0104m"
	BrightMagentaBackground       = "\x1b[0105m"
	BrightCyanBackground          = "\x1b[0106m"
	BrightWhiteBackground         = "\x1b[0107m"
)

const underlineChar = "="

var attrColour = color.New(color.FgBlue, color.Bold)
var headingColour = color.New(color.FgGreen)
var warningColour = color.New(color.FgYellow, color.Bold)

func formatAttrString(attr string, value string) string {
	return fmt.Sprintf("%s: %s", attrColour.Sprint(attr), value)
}

func formatAttrInt(attr string, value int) string {
	return fmt.Sprintf("%s: %d", attrColour.Sprint(attr), value)
}

func printHeading(heading string) {
	fmt.Printf("%s\n", headingColour.Sprint(heading))
	underlineStr := strings.Repeat(underlineChar, len(heading))
	fmt.Printf("%s\n", headingColour.Sprint(underlineStr))
}

func writeHeading(writer io.Writer, header []string) {
	writeRow(writer, printUniformRow(GreenText, header))
	writeRow(writer, printUniformRow(GreenText, anonymizeRow(header))) // header separator
}

func printWarningString(warning string) {
	fmt.Printf("warning: %s\n", warningColour.Sprint(warning))
}

// Color implements the Stringer interface for interoperability with string
func (c *Color) String() string {
	return fmt.Sprintf("%v", c)
}

func printColStr(color Color, value string) string {
	return fmt.Sprintf("%v%v%v", color, value, Reset)
}

func printColRow(colors []Color, row []string) []string {
	paintedRow := make([]string, len(row))
	for i, v := range row {
		paintedRow[i] = printColStr(colors[i], v)
	}
	return paintedRow
}

type priceField float32
type percentField float32

var priceFmt = "%8.8f"
var percentFmt = "%3.2f"

func formatColRow(cols ...interface{}) []string {
	formattedRow := make([]string, len(cols))
	for i, col := range cols {

		if str, ok := col.(string); ok {
			// is string
			formattedRow[i] = printColStr(WhiteText, str)
		} else {

			// Set default
			color := Color(BrightGreenText)
			colFmt := "%s"

			num := float32(999.999)

			// override by specific type
			switch col.(type) {
			case priceField:
				colFmt = priceFmt
				priceNum := col.(priceField)
				if priceNum < 0 {
					color = BrightRedText
				}
				num = float32(priceNum)
				formattedRow[i] = printColStr(color, fmt.Sprintf(colFmt, num))
			case percentField:
				colFmt = percentFmt
				pctNum := col.(percentField)
				if pctNum < 0 {
					color = BrightRedText
				}
				num = float32(pctNum)
				formattedRow[i] = printColStr(color, fmt.Sprintf(colFmt, num))
			default:
				formattedRow[i] = printColStr(BrightYellowText, fmt.Sprintf("Unexpected: %#v", col))

			}
		}
	}
	return formattedRow
}

func printUniformRow(color Color, row []string) []string {
	colors := make([]Color, len(row))
	for i, _ := range colors {
		colors[i] = color
	}
	return printColRow(colors, row)
}

func anonymizeRow(row []string) []string {
	anonRow := make([]string, len(row))
	for i, v := range row {
		anonRow[i] = strings.Repeat("-", len(v))
	}
	return anonRow
}

func writeRow(writer io.Writer, line []string) {
	fmt.Fprintln(writer, strings.Join(line, "\t"))
}

func printErr(err error) string {
	return printColStr(BrightRedBackground+BrightWhiteText, "ERROR: "+err.Error())
}

func formatProtoTimestamp(ts *google_protobuf.Timestamp) string {
	if ts == nil {
		return ""
	}
	if tt, err := tspb.Timestamp(ts); err != nil {
		return fmt.Sprintf("Invalid timestamp: %s", err)
	} else {
		return tt.Format(DATE_FORMAT)
	}
}
