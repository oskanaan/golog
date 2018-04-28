package logreader

import (
	"bytes"
	"io"
	"os"
	"strings"
	"bufio"
)

//Reads a maximum of "capacity" number of lines starting from the offset position
//Returns a slice containing a maximum of "capacity" entries and the current position.
func readLinesStartingFromPosition(input io.ReadSeeker, capacity int, start int) ([]string, int, error) {
	rows := make([]string, 0)
	if _, err := input.Seek(int64(start), 0); err != nil {
		return rows, -1, err
	}

	r := bufio.NewReader(input)
	pos := int64(start)
	for i := 0; i < capacity; i++ {
		data, err := r.ReadBytes('\n')
		pos += int64(len(data))
		if err == nil || err == io.EOF {
			if len(data) > 1 && data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
			//Only a new line character
			if len(data) == 1 && data[len(data)-1] == '\n' {
				data = []byte{' '}
			}

			if len(data) > 0 && data[len(data)-1] == '\r' {
				data = data[:len(data)-1]
			}
			rows = append(rows, string(data))
		}
		if err != nil {
			if err != io.EOF {
				return rows, int(pos), err
			}
			break
		}
	}
	return rows, int(pos), nil
}

//A convenience method for tailing a file
//Returns a slice containing the retrieved rows and the new offset
func tail(file *os.File, capacity int, endOffset int) ([]string, int, error) {
	tailStartPosition := tailStartPosition(file, capacity, endOffset)
	return readLinesStartingFromPosition(file, capacity, tailStartPosition)
}

//A convenience method to head a file
//Returns a slice containing the retrieved rows and the new offset
func head(file *os.File, capacity int, offset int) ([]string, int, error) {
	return readLinesStartingFromPosition(file, capacity, offset)
}

func nextLine(file *os.File, offset int) (string, int, error) {
	data, offset, err := head(file, 1, offset)
	return data[0], offset, err
}

func lastLine(file *os.File, capacity int, offset int) (string, int, error) {
	data, offset, err := tail(file, capacity, offset)
	return data[capacity-1], offset, err
}

func tailStartPosition(file *os.File, capacity int, endOffset int) int {
	bufferSize := 128
	fileInfo, _ := file.Stat()
	size := endOffset
	if endOffset == -1 {
		size = int(fileInfo.Size())
	}
	buf := make([]byte, bufferSize)
	lineSep := []byte{'\n'}

	linesFromEnd := 0
	newOffset := size

	for linesFromEnd < capacity {
		newOffset -= bufferSize
		if newOffset < 0 {
			bufDifference := newOffset
			newOffset = 0
			if len(buf) <= -1 * bufDifference {
				buf = buf[0:0]
			} else {
				buf = buf[:len(buf) + bufDifference]
			}
		}
		if len(buf) == 0 {
			break
		}
		file.Seek(int64(newOffset), 0)
		_, err := file.Read(buf)

		if err != nil {
			break
		}
		if bytes.Count(buf, lineSep) >= 0 {
			linesFromEnd += bytes.Count(buf, lineSep)
			//Read too much, go to the correct offset.
			for linesFromEnd > capacity {
				newLineIndex := bytes.Index(buf, lineSep)
				if newLineIndex == -1 {
					break
				} else {
					//Get past the new line character
					newLineIndex ++
				}
				buf = buf[newLineIndex:]
				newOffset += newLineIndex
				linesFromEnd--

			}
		}
	}

	if newOffset < 0 {
		newOffset = 0
	} else {
		if newOffset > 0 {
			file.Seek(int64(newOffset-1), 0)
			previousChar := make([]byte, 1)
			file.Read(previousChar)

			if previousChar[0] != '\n' {
				//Are we at the beginning of the line? if not, get the offset of the start of the first line within 512 characters.
				lineBufSize := 512
				startOffset := newOffset - lineBufSize - 1
				if startOffset < 0 {
					lineBufSize += startOffset
					startOffset = 0
				}

				if lineBufSize > 0 {
					buf := make([]byte, lineBufSize)
					file.Seek(int64(startOffset), 0)
					file.Read(buf)

					if buf[lineBufSize-1] != '\n' {
						lastNewLineIndex := bytes.LastIndex(buf, lineSep)
						newOffset -= len(string(buf[lastNewLineIndex+1:])) + 1
					}
				}
			}
		}
	}

	return newOffset
}

func readLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	scanError := io.EOF

	for {
		for sc.Scan() {
			lastLine++
			if lastLine == lineNum {
				// you can return sc.Bytes() if you need output in []bytes
				return sc.Text(), lastLine, sc.Err()
			}
		}
		//In case the last scan caused any issues, such as ErrorTooLong.
		if sc.Err() != nil {
			scanError = sc.Err()
		}

		if scanError == io.EOF {
			break
		}
	}

	return line, lastLine, scanError
}

//Reads all lines related to a stack trace starting from the specified index "lineNum"
//Returns a string representing the stack trace
func stackTrace(r *os.File, offset int, delim string) (stackTrace string) {
	tailStart := tailStartPosition(r, 20, offset)
	data, newOffset, _ := head(r, 20, tailStart)
	linesRead := 0

	for {
		for _, line := range data {
			if !strings.Contains(line, delim) {
				stackTrace = stackTrace + "\n" + line
			} else {
				return stackTrace
			}

			//A maximum of a 100 lines to be displayed
			//Todo: Display the stack trace in a scrollable view.
			if linesRead > 100 {
				return stackTrace
			}

			linesRead++
		}
		data, newOffset, _ = head(r, 20, newOffset)

		if len(data) == 0 {
			break
		}
	}


	return stackTrace
}
