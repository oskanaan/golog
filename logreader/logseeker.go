package logreader

import (
	"os"
	"io"
	"bytes"
	"bufio"
)

const readBufSize  = 5120

const(
	beginning = iota
	middle
	end
	ignore
)

//Reads a maximum of "capacity" number of lines starting from the offset position
//from the end of the file.
//Returns a slice containing a maximum of "capacity" entries and the current position.
func readFileFromEnd(file *os.File, capacity, lineNumber int) ([]string, int) {
	rows := make([] string, 0)
	linesRead := 0
	currentPosition := lineNumber
	for linesRead < capacity {
		//If first iteration and the cursor will move before the file beginning, set the current position to 1 position after the capacity
		if currentPosition <= capacity && linesRead == 0 {
			currentPosition = capacity + 1
		}
		line, _ := getPreviousLine(file, currentPosition)
		currentPosition--

		if line == "" {
			continue
		}

		tempRows := make([] string, 0)
		tempRows = append(tempRows, line)
		rows = append(tempRows, rows ...)
		linesRead ++
	}

	return rows, currentPosition
}

//A convenience method for tailing a file
func tail(file *os.File, capacity int) ([]string, int) {
	fileLineCount, _ := getLineCount(file)
	return readFileFromEnd(file, capacity, fileLineCount+1)
}

func getPreviousLine(file *os.File, currentPosition int) (string, int) {
	previousLineNumber := currentPosition - 1
	if previousLineNumber < 0 {
		return "", beginning
	}

	file.Seek(0, 0)
	line, _, _ := readLine(file, previousLineNumber)

	return line, middle
}

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

//func getPreviousLine(file *os.File, offset, bufferDelta int) (string, int) {
//	var read string
//	bufferSize := readBufSize
//	newOffset := offset
//	fileInfo, _ := file.Stat()
//	logSize := fileInfo.Size()
//	isAtBegnning := false
//
//	for {
//		seekOffset := int64(newOffset - (bufferSize - bufferDelta))
//		if -1 * seekOffset > logSize {
//			bufferSize -= -1 * int(logSize + seekOffset)
//			if bufferSize < bufferDelta {
//				bufferSize = bufferDelta - 1
//			}
//			seekOffset = -1 * logSize
//			isAtBegnning = true
//		}
//		ret, err := file.Seek(seekOffset,2)
//		if err != nil {
//			fmt.Println(err, ret)
//		}
//
//		newOffset = int(ret)
//		buf := make([]byte, bufferSize )
//		file.Read(buf)
//		read = string(buf) + read
//
//		if isAtBegnning {
//			lines := strings.Split(read, "\n")
//			return lines[len(lines)-1], beginning
//		}
//
//		if strings.Contains(read, "\n"){
//			lines := strings.Split(read, "\n")
//			return lines[len(lines)-1], middle
//		}
//
//	}
//
//	return "",ignore
//}