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
		rows = append(rows, parseLine(line, l.config.Delim))
	}

	l.currentOffset = offset
	return &rows
}

//Reads the first N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Head() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data ,offset := head(file, l.config.Capacity)
	rows := [][] string {}
	for _, line := range data {
		rows = append(rows, parseLine(line, l.config.Delim))
	}

	l.currentOffset = offset
	return &rows
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last page
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageUp() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Delim, l.config.Capacity, l.currentOffset)
	l.currentOffset = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the first line after the current page
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageDown() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Delim, l.config.Capacity, l.currentOffset + (l.config.Capacity * 2))
	l.currentOffset = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last line
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Up() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Delim, l.config.Capacity, l.currentOffset + l.config.Capacity - 1)
	l.currentOffset = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the current line + 1
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Down() *[][]string {
	file ,_ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Delim, l.config.Capacity, l.currentOffset + l.config.Capacity + 1 )
	l.currentOffset = offset
	return data
}

func countLines(lines string) int {
	return len(strings.Split(lines, "\n"))
}

//Parses a single log file line using the "delim" character to separate the columns
//Returns a slice of strings containing the row data
func parseLine(line string, delim string) []string {
	if line == "" {
		return []string{}
	}

	columns := strings.Split(line, delim)
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

//Sets the number of rows to display capacity
func (l *LogReader) SetCapacity(capacity int) {
	l.config.Capacity = capacity
}

func (l *LogReader) Message(lineNum int) string{
	file ,_ := os.Open(l.input)
	defer file.Close()

	message, _, _ := readLine(file, lineNum + l.currentOffset - 1)
	if !strings.Contains(message, l.config.Delim) {
		file.Seek(0,0)
		message = stackTrace(file, lineNum + l.currentOffset - 1, l.config.Delim)
	}

	return message
}

//Reads N (N=capacity) lines starting from the offset
//Returns a two dimensional array containing the parsed columns and the new offset
func readLogFileFromOffset(file *os.File, delim string, capacity int, offset int) (*[][]string, int) {
	data ,offset := readFileFromEnd(file, capacity, offset )
	rows := [][] string {}
	if len(data) == 0 {
		return &rows, 0
	}
	for _, line := range data {
		rows = append(rows, parseLine(line, delim))
	}

	return &rows, offset
}
