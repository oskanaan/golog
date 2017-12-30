package logdisplay

import (
	"fmt"
	"strings"
)

const err = "Error"
const warn = "Warn"
const info = "Info"
const trace = "Trace"
const debug = "Debug"

//prepends/appends the header color to the log headers
//Returns a color coded string
func colorizeHeader(header string) string {
	return fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2, 4), header, "\033[0m")
}

//prepends/appends the color of the log entry based on the severity
//Returns a color coded string
//Todo: make the colors and severity values configurable
func colorizeLogEntry(logEntry string) string {
	var pre string
	if caseInsensitiveContains(logEntry, err) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 1, 1)
	} else if caseInsensitiveContains(logEntry, warn) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 3, 1)
	} else if caseInsensitiveContains(logEntry, debug) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 2, 1)
	} else if caseInsensitiveContains(logEntry, info) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 6, 5)
	} else if caseInsensitiveContains(logEntry, trace) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 0, 1)
	} else {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 0, 1)
	}
	return fmt.Sprint(pre, logEntry, "\033[0m")
}

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}