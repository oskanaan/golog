package logreader

import (
	"strings"
	"os"
	"bufio"
)

type LogReader struct {
	input  string
	config Config
}

//Config is used to configure the behaviour of the log reader
type Config struct {
	Delim string
	Headers []string
	ColumnSizes []int
	Capacity int
}

func NewLogReader(input string, config Config ) LogReader {
	var l LogReader
	l.input = input
	l.config = config

	return l
}

//Reads the last N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l LogReader) Tail() [][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rows := [][] string {}

	for scanner.Scan() {
		if len(rows) >= l.config.Capacity {
			rows = rows[1:]
		}
		rows = append(rows, l.parseLine(scanner.Text()))
	}
	return rows
}

//Gets a slice of strings representing the headers of the log
func (l LogReader) GetHeaders() [] string{
	return l.config.Headers
}

//Gets a slice of strings representing the headers of the log
func (l LogReader) GetColumnSizes() [] int{
	return l.config.ColumnSizes
}

//Parses a single log file line using the "delim" character to separate the columns
//Returns a slice of strings containing the row data
func (l LogReader) parseLine(line string) []string {
	if line == "" {
		return []string{}
	}

	columns := strings.Split(line, l.config.Delim)
	var columnValues []string

	for _, col := range columns {
		columnValues = append(columnValues, col)
	}

	return columnValues
}