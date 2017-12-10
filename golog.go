package main

import (
	"flag"
	"github.com/oskanaan/golog/logreader"
	"strings"
	"strconv"
	"syscall"
	"fmt"
	"github.com/oskanaan/golog/logdisplay"
)

func main() {
	//Read command line arguments
	file := flag.String("file", "test.log", "Log file to view")
	seperator := flag.String("seperator", "~", "Log column seperator")
	headersString := flag.String("headers", "", "Comma seperated log columns header labels")
	columnSizesString := flag.String("column-sizes", "", "Comma seperated list of columns sizes in characters")
	capacityString := flag.String("capacity", "", "Number of lines to display at a time")
	flag.Parse()

	headers := strings.Split(*headersString, ",")
	columnSizes := func()[]int {
		var sizes []int
		for _,val := range strings.Split(*columnSizesString, ","){
			i,err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				fmt.Printf("Could not parse column size %s to number\n", val)
				syscall.Exit(1)
			}
			sizes = append(sizes, int(i))
		}
		return sizes
	}()

	capacity, err := strconv.Atoi(*capacityString)
	if err != nil {
		fmt.Errorf("terminal-capacity should be a integer")
		syscall.Exit(1)
	}

	logReader := logreader.NewLogReader(*file, logreader.Config{*seperator, headers,  columnSizes, capacity})
	logDisplay := logdisplay.NewLogDisplay(logReader)
	logDisplay.Display()
}
