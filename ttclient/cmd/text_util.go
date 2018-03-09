package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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
func printWarningString(warning string) {
	fmt.Printf("warning: %s\n", warningColour.Sprint(warning))
}
