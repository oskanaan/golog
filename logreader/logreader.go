package logreader

import (
	"strings"
	"os"
)

type LogReader struct {
	input  string
	config Config
	line int
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
func (l LogReader) Tail() [][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	rows := [][] string {}
	linesRead := 0
	currentPosition := -2000
	for linesRead < l.config.Capacity {
		file.Seek(int64(currentPosition),2)
		buf := make([]byte, 2000 )
		file.Read(buf)
		lines := strings.Split(string(buf), "\n")
		for i:=len(lines)-1 ; i>0 && linesRead < l.config.Capacity ; i-- {
			tempRows := [][] string {}
			tempRows = append(tempRows, l.parseLine(lines[i]))
			rows = append(tempRows, rows ...)
			linesRead ++
			currentPosition -= len(lines[i]) + 1
		}
	}

	return rows
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