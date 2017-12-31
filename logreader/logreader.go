package logreader

import (
	"os"
	"strings"
)

type LogReaderConfig struct {
	LogFile string `yaml:"file"`
	Seperator string `yaml:"seperator"`
	Headers []struct {
		Header string `yaml:"header"`
		Size int `yaml:"size"`
	}
}

type LogReader struct {
	input         string
	config        LogReaderConfig
	currentOffset int
	Capacity      int
}

func NewLogReader(input string, config LogReaderConfig) LogReader {
	var l LogReader
	l.input = input
	l.config = config

	return l
}

//Reads the last N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Tail() *[][]string {
	file, _ := os.Open(l.input)
	defer file.Close()

	data, offset := tail(file, l.Capacity)
	rows := [][]string{}
	for _, line := range data {
		rows = append(rows, parseLine(line, l.config.Seperator))
	}

	l.currentOffset = offset
	return &rows
}

//Reads the first N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Head() *[][]string {
	file, _ := os.Open(l.input)
	defer file.Close()

	data, offset := head(file, l.Capacity)
	rows := [][]string{}
	for _, line := range data {
		rows = append(rows, parseLine(line, l.config.Seperator))
	}

	l.currentOffset = offset
	return &rows
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last page
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageUp() *[][]string {
	file, _ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Seperator, l.Capacity, l.currentOffset)
	l.currentOffset = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the first line after the current page
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageDown() *[][]string {
	file, _ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Seperator, l.Capacity, l.currentOffset+(l.Capacity*2))
	l.currentOffset = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last line
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Up() *[][]string {
	file, _ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Seperator, l.Capacity, l.currentOffset+l.Capacity-1)
	l.currentOffset = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the current line + 1
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Down() *[][]string {
	file, _ := os.Open(l.input)
	defer file.Close()

	data, offset := readLogFileFromOffset(file, l.config.Seperator, l.Capacity, l.currentOffset+l.Capacity+1)
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
func (l LogReader) GetHeaders() []string {
	headers := make([]string, len(l.config.Headers))
	for index, header := range l.config.Headers {
		headers[index] = header.Header
	}
	return headers
}

//Gets a slice of strings representing the headers of the log
func (l LogReader) GetColumnSizes() []int {
	sizes := make([]int, len(l.config.Headers))
	for index, header := range l.config.Headers {
		sizes[index] = header.Size
	}
	return sizes
}

//Sets the number of rows to display capacity
func (l *LogReader) SetCapacity(capacity int) {
	l.Capacity = capacity
}

func (l *LogReader) Message(lineNum int) string {
	file, _ := os.Open(l.input)
	defer file.Close()

	message, _, _ := readLine(file, lineNum+l.currentOffset-1)
	if !strings.Contains(message, l.config.Seperator) {
		file.Seek(0, 0)
		message = stackTrace(file, lineNum+l.currentOffset-1, l.config.Seperator)
	}

	return message
}

//Reads N (N=capacity) lines starting from the offset
//Returns a two dimensional array containing the parsed columns and the new offset
func readLogFileFromOffset(file *os.File, delim string, capacity int, offset int) (*[][]string, int) {
	data, offset := readLinesStartingFromPosition(file, capacity, offset)
	rows := [][]string{}
	if len(data) == 0 {
		return &rows, 0
	}
	for _, line := range data {
		rows = append(rows, parseLine(line, delim))
	}

	return &rows, offset
}
