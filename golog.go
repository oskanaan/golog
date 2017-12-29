package main

import (
	"flag"
	"fmt"
	"github.com/oskanaan/golog/logdisplay"
	"github.com/oskanaan/golog/logreader"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	//Read command line arguments
	file := flag.String("file", "test.log", "Log file to view")
	seperator := flag.String("seperator", "~", "Log column seperator")
	headersString := flag.String("headers", "", "Comma seperated log columns header labels")
	columnSizesString := flag.String("column-sizes", "", "Comma seperated list of columns sizes in characters")
	severityColumn := flag.String("severity", "Severity", "The column which determines the severity of the log line")
	flag.Parse()

	headers := strings.Split(*headersString, ",")
	columnSizes := func() []int {
		var sizes []int
		for _, val := range strings.Split(*columnSizesString, ",") {
			i, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				fmt.Printf("Could not parse column size %s to number\n", val)
				syscall.Exit(1)
			}
			sizes = append(sizes, int(i))
		}
		return sizes
	}()

	logReader := logreader.NewLogReader(*file, logreader.Config{*seperator, headers, columnSizes, 10, *severityColumn})
	logDisplay := logdisplay.NewLogDisplay(&logReader)
	logDisplay.DisplayUI()
}
