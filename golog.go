package main

import (
	"flag"
	"bufio"
	"os"
	"github.com/oskanaan/golog/logreader"
	"strings"
	"strconv"
	"syscall"
	"fmt"
)

func main() {
	file := flag.String("file", "test.log", "Log file to view")
	seperator := flag.String("seperator", "~", "Log column seperator")
	headersString := flag.String("headers", "", "Comma seperated log columns header labels")
	columnSizesString := flag.String("column-sizes", "", "Comma seperated list of columns sizes in characters")
	flag.Parse()

	headers := strings.Split(*headersString, ",")
	columnSizes := func()[]int {
		var sizes []int
		for _,val := range strings.Split(*columnSizesString, ","){
			i,err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				fmt.Errorf("Could not parse column size %s to number", val)
				syscall.Exit(1)
			}
			sizes = append(sizes, int(i))
		}
		return sizes
	}()
	out := bufio.NewWriter(os.Stdout)
	logReader := logreader.NewLogReader(*file, *out, logreader.Config{*seperator, headers,  columnSizes})
	logReader.StartReading()
}
