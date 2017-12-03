package logreader

import (
	"bufio"
	"strings"
	"github.com/hpcloud/tail"
	"fmt"
	"strconv"
)

type LogReader struct {
	input  string
	output bufio.Writer
	config Config
}

//Config is used to configure the behaviour of the log reader
type Config struct {
	Delim string
	Headers []string
	ColumnSizes []int
}

func NewLogReader(input string, output bufio.Writer, config Config ) LogReader {
	var l LogReader
	l.input = input
	l.output = output
	l.config = config

	return l
}

//Starts tailing the log file
func (l LogReader) StartReading() {
	t, _ := tail.TailFile(l.input, tail.Config{Follow: true, ReOpen: true, Poll: true})
	l.writeHeader()

	for line := range t.Lines {
		l.parseLine(line.Text)
	}
}

//Writes the header values in the "headers" slice and rpad them using the column size as the corresponding index
func (l LogReader) writeHeader() {
	var format string
	var output string

	for i:=0 ; i<len(l.config.Headers) ; i++ {
		format = output
		//rpad columns using the configured column size
		format += "%-"+strconv.Itoa(l.config.ColumnSizes[i])+"s"
		output = fmt.Sprintf(format, l.config.Headers[i])
	}
	l.output.WriteString(output)
	l.output.WriteString("\n")
	l.output.Flush()
}

//Parses a single log file line using the "delim" character to separate the columns
//Returns a boolean indicating that the parse was successful or not
func (l LogReader) parseLine(line string) bool {
	if line == "" {
		return false
	}

	columns := strings.Split(line, l.config.Delim)
	var columnValues []string
	var format string
	var output string

	for index, col := range columns {
		columnValues = append(columnValues, col)
		format = output
		//rpad columns using the configured column size
		format += "%-"+strconv.Itoa(l.config.ColumnSizes[index])+"s"
		output = fmt.Sprintf(format, col)
	}

	//Display the formatted line on the terminal 
	l.output.WriteString(output)
	l.output.WriteString("\n")
	l.output.Flush()

	return true
}
