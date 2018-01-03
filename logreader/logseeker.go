package logreader

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

const readBufSize = 5120

//Reads a maximum of "capacity" number of lines starting from the offset position
//Returns a slice containing a maximum of "capacity" entries and the current position.
func readLinesStartingFromPosition(file *os.File, capacity, lineNumber int) ([]string, int) {
	rows := make([]string, 0)
	linesRead := 0
	currentPosition := lineNumber
	for linesRead < capacity {
		//If first iteration and the cursor will move before the file beginning, set the current position to 1 position after the capacity
		if currentPosition <= capacity && linesRead == 0 {
			currentPosition = capacity + 1
		}
		line := getPreviousLine(file, currentPosition)
		currentPosition--

		if line == "" {
			continue
		}

		tempRows := make([]string, 0)
		tempRows = append(tempRows, line)
		rows = append(tempRows, rows...)
		linesRead++
	}

	return rows, currentPosition
}

//A convenience method for tailing a file
//Returns a slice containing the retrieved rows and the new offset
func tail(file *os.File, capacity int) ([]string, int) {
	fileLineCount, _ := getLineCount(file)
	return readLinesStartingFromPosition(file, capacity, fileLineCount+1)
}

//A convenience method to head a file
//Returns a slice containing the retrieved rows and the new offset
func head(file *os.File, capacity int) ([]string, int) {
	return readLinesStartingFromPosition(file, capacity, 0)
}

//Retrieves the previous line starting from the current position
//Returns a the previous line and the new offset
func getPreviousLine(file *os.File, currentPosition int) string {
	previousLineNumber := currentPosition - 1
	if previousLineNumber < 0 {
		return ""
	}

	file.Seek(0, 0)
	line, _, _ := readLine(file, previousLineNumber)

	return line
}

//Counts the total number of lines within a file
//Returns the total count and an error if a non EOF error is encountered
func getLineCount(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

//Reads a specific line from a file
//Returns the line at the specified index, the index of the last line read, and an error if any error was encountered
func readLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

//Reads all lines related to a stack trace starting from the specified index "lineNum"
//Returns a string representing the stack trace
func stackTrace(r *os.File, lineNum int, delim string) (stackTrace string) {
	index := lineNum
	linesRead := 1
	line := ""

	for !strings.Contains(line, delim) {
		r.Seek(0, io.SeekStart)
		stackTrace = stackTrace + "\n" + line

		//A maximum of a 100 lines to be displayed
		//Todo: Display the stack trace in a scrollable view.
		if linesRead > 100 {
			break
		}

		line, _, _ = readLine(r, index)
		index++
		linesRead++

		//if s, ok := r.(io.Seeker); ok {
		//	s.Seek(int64(len(line))+1, io.SeekCurrent) // seek relative to current file pointer
		//}
		//line, _, _ = readLine(r, linesRead)
		//fmt.Println(line, strings.Contains(line, delim), index, linesRead)

	}

	return stackTrace
}
