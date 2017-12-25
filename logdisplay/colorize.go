//colorize has some utility functions for adding color codes to the log lines
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
func colorizeHeader(header string) string {
	return fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2,4 ), header, "\033[0m")
}

//prepends/appends the color of the log entry based on the severity
func colorizeLogEntry(logEntry, severity string) string {
	var pre string
	if strings.EqualFold(severity, err) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 1, 1 )
	} else if strings.EqualFold(severity, warn) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 3, 1 )
	} else if strings.EqualFold(severity, debug) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 2, 1 )
	} else if strings.EqualFold(severity, info) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 6, 5 )
	} else if strings.EqualFold(severity, trace) {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 0, 1 )
	} else {
		pre = fmt.Sprintf("\033[3%d;%d;1m", 0, 1 )
	}
	return fmt.Sprint(pre, logEntry, "\033[0m")
}