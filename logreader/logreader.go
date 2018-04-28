package logreader

import (
	"os"
	"strings"
	"bufio"
	"fmt"
)

type LogReaderConfig struct {
	AutoDetectFiles     string `yaml:"autoDetectFiles"`
	Files               []LogFile `yaml:files`
	Seperator           string `yaml:"seperator"`
	Headers             []Header
}

type Header struct {
	Header string `yaml:"header"`
	Size int `yaml:"size"`
}

type LogFile struct {
	LogFile string `yaml:"file"`
	Name    string `yaml:"name"`
}

type LogReader struct {
	FileIndex     int
	config        LogReaderConfig
	currentOffset []int
	Capacity      int
	currentLoadedPage *[][]string
	previousReadFileInfo os.FileInfo
}

func NewLogReader(config LogReaderConfig) LogReader {
	var l LogReader
	l.config = config
	l.currentOffset = make([]int, len(l.config.Files))

	return l
}

//Reads the log file from the current offset
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Refresh() *[][]string {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	data, offset := readLogFileFromOffsetUp(file, l.config.Seperator, l.Capacity, l.currentOffset[l.FileIndex] + l.Capacity)
	l.currentOffset[l.FileIndex] = offset
	return data
}


//Reads the last N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Tail() *[][]string {
	fileInfo, err := os.Stat(l.config.Files[l.FileIndex].LogFile)
	if err != nil {
		return &[][]string{}
	}
	//Check if any changes happened to the file, otherwise return last page
	if l.previousReadFileInfo != nil &&
		fileInfo.Size() == l.previousReadFileInfo.Size() &&
		fileInfo.ModTime() == l.previousReadFileInfo.ModTime() &&
		l.currentLoadedPage != nil {

		l.currentOffset[l.FileIndex] = int(fileInfo.Size())
		return l.currentLoadedPage
	}

	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	data, _, _ := tail(file, l.Capacity, -1)
	rows := [][]string{}
	for _, line := range data {
		rows = append(rows, parseLine(line, l.config.Seperator))
	}

	l.currentOffset[l.FileIndex] = int(fileInfo.Size())
	l.currentLoadedPage = &rows
	l.previousReadFileInfo = fileInfo
	return &rows
}

//Reads the first N lines where N=The capacity configuration value
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Head() *[][]string {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	data, offset, _ := head(file, l.Capacity, 0)
	rows := [][]string{}
	for _, line := range data {
		rows = append(rows, parseLine(line, l.config.Seperator))
	}
	l.currentOffset[l.FileIndex] = offset
	return &rows
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last page
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageUp() *[][]string {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	data, offset := readLogFileFromOffsetUp(file, l.config.Seperator, l.Capacity, l.currentOffset[l.FileIndex])
	l.currentOffset[l.FileIndex] = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the first line after the current page
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) PageDown() *[][]string {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	data, offset := readLogFileFromOffsetDown(file, l.config.Seperator, l.Capacity, l.currentOffset[l.FileIndex])
	l.currentOffset[l.FileIndex] = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the current offset excluding the last line
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Up() *[][]string {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	lastLine, _, _ := lastLine(file, l.Capacity, l.currentOffset[l.FileIndex])
	lastPosition := tailStartPosition(file, l.Capacity, l.currentOffset[l.FileIndex]-len(lastLine))
	if lastPosition < 0 {
		lastPosition = 0
	}

	data, offset := readLogFileFromOffsetDown(file, l.config.Seperator, l.Capacity, lastPosition )
	l.currentOffset[l.FileIndex] = offset
	return data
}

//Reads the last N lines where N=The capacity configuration value starting from the current line + 1
//Returns a two dimensional slice containing the parsed rows
func (l *LogReader) Down() *[][]string {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}
	}
	defer file.Close()

	nextLine, _, _ := nextLine(file, l.currentOffset[l.FileIndex])
	data, offset := readLogFileFromOffsetDown(file, l.config.Seperator, l.Capacity, tailStartPosition(file, l.Capacity, l.currentOffset[l.FileIndex]+len(nextLine)+1))
	l.currentOffset[l.FileIndex] = offset
	return data
}

//Search the log file for a search term
//Returns a two dimensional slice containing the parsed rows for the location containing the search term and the location within the result
func (l *LogReader) Search(searchTerm string, currentLocation int) (*[][]string, int) {
	file, err := l.openLogFile()
	if err != nil {
		return &[][]string{}, 0
	}
	defer file.Close()

	searchOffset := l.currentOffset[l.FileIndex] + currentLocation + 1
	location := searchFileForTerm(file, searchTerm, searchOffset)
	if location == -1 {
		file.Seek(0, 0)
		location = searchFileForTerm(file, searchTerm, 0)
	}
	startFromLocation := location + l.Capacity
	data, offset := readLogFileFromOffsetUp(file, l.config.Seperator, l.Capacity, startFromLocation)
	l.currentOffset[l.FileIndex] = offset
	resultLocationInCurrentPage := offset - location

	return data, resultLocationInCurrentPage
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

//Sets the number of rows to display capacity
func (l *LogReader) CurrentOffset() int {
	return l.currentOffset[l.FileIndex]
}

func (l *LogReader) Message(lineNum int) string {
	file, err := l.openLogFile()
	if err != nil {
		return ""
	}

	defer file.Close()

	data, _, _ := tail(file, l.Capacity, l.currentOffset[l.FileIndex])

	message := data[lineNum - 1]
	fmt.Println(message)
	if !strings.Contains(message, l.config.Seperator) {
		message = stackTrace(file, l.currentOffset[l.FileIndex], l.config.Seperator)
	}

	return message
}

func (l *LogReader) Progress() int {
	file, err := l.openLogFile()
	if err != nil {
		return -1
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	percentage := int(float32(l.currentOffset[l.FileIndex])/float32(fileInfo.Size()) * 100)
	return percentage
}

//Reads N (N=capacity) lines starting from the offset
//Returns a two dimensional array containing the parsed columns and the new offset
func readLogFileFromOffsetUp(file *os.File, delim string, capacity int, offset int) (*[][]string, int) {
	data, newOffset, _ := tail(file, capacity, tailStartPosition(file, capacity, offset))
	rows := [][]string{}
	if len(data) == 0 {
		return &rows, 0
	}
	for _, line := range data {
		rows = append(rows, parseLine(line, delim))
	}

	return &rows, int(newOffset)
}

//Reads N (N=capacity) lines starting from the offset
//Returns a two dimensional array containing the parsed columns and the new offset
func readLogFileFromOffsetDown(file *os.File, delim string, capacity int, offset int) (*[][]string, int) {
	fileInfo, _ := file.Stat()
	data, newOffset, _ := head(file, capacity, offset)
	if newOffset >= int(fileInfo.Size()) {
		data, newOffset, _ = tail(file, capacity, -1)
	}

	rows := [][]string{}
	if len(data) == 0 {
		return &rows, 0
	}
	for _, line := range data {
		rows = append(rows, parseLine(line, delim))
	}

	return &rows, int(newOffset)
}

func searchFileForTerm(f *os.File, searchTerm string, currentLocation int) int{
	scanner := bufio.NewScanner(f)
	line := 1
	for scanner.Scan() {
		if line > currentLocation && strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(searchTerm)) {
			return line
		}
		line++
	}

	return -1
}

func (l LogReader) openLogFile() (*os.File, error) {
	return os.Open(l.config.Files[l.FileIndex].LogFile)
}