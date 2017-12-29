package logreader

import (
	"strings"
	"os"
)

type LogReader struct {
	input         string
	config        Config
	currentOffset int
}

//Config is used to configure the behaviour of the log reader
type Config struct {
	Delim string
	Headers []string
	ColumnSizes []int
	Capacity int
	SeverityColumn string
}

func NewLogReader(input string, config Config ) LogReader {
	var l LogReader
	l.input = input
	l.config = config

	return l
}

//Reads the last N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Tail() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data ,offset := tail(file, l.config.Capacity)
	rows := [][] string {}
	for _, line := range data {
		rows = append(rows, l.parseLine(line))
	}

	l.currentOffset = offset
	return &rows
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last line
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageUp() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data ,offset := readFileFromEnd(file, l.config.Capacity, l.currentOffset )
	rows := [][] string {}
	if len(data) == 0 {
		return &rows
	}
	for _, line := range data {
		rows = append(rows, l.parseLine(line))
	}

	l.currentOffset = offset
	return &rows
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset including the line after the last
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageDown() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data ,offset := readFileFromEnd(file, l.config.Capacity, l.currentOffset + (l.config.Capacity * 2))
	rows := [][] string {}
	for _, line := range data {
		rows = append(rows, l.parseLine(line))
	}

	l.currentOffset = offset
	return &rows
}


func countLines(lines string) int {
	return len(strings.Split(lines, "\n"))
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

//Gets a slice of strings representing the headers of the log
func (l LogReader) GetHeaders() [] string{
	return l.config.Headers
}

//Gets a slice of strings representing the headers of the log
func (l LogReader) GetColumnSizes() [] int{
	return l.config.ColumnSizes
}

//Gets a the severity column name
func (l LogReader) GetSeverityColumnName() string{
	return l.config.SeverityColumn
}