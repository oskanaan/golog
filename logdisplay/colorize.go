package logdisplay

import (
	"fmt"
	"regexp"
)

var severity = []string {"ERROR", "WARN", "INFO", "TRACE", "DEBUG"}
var colorCodes = [][]interface{} {{1,1}, {3,1}, {2,1}, {6,5}, {0,1}}

//prepends/appends the header color to the log headers
//Returns a color coded string
func colorizeHeader(header string) string {
	return fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2, 4), header, "\033[0m")
}

//prepends/appends the color of the log entry based on the severity
//Returns a color coded string
//Todo: make the colors and severity values configurable
func colorizeLogEntry(logEntry string) string {
	colorCode := severityColorCode(logEntry)
	pre := fmt.Sprintf("\033[3%d;%d;1m", colorCode...)
	return fmt.Sprint(pre, logEntry, "\033[0m")
}

func severityColorCode(entry string) []interface{} {
	for index, sev := range severity {
		if severityMatch(entry, sev) {
			return colorCodes[index]
		}
	}

	return colorCodes[len(colorCodes) - 1]
}

func severityMatch(s, substr string) bool {
	match, _ := regexp.MatchString(fmt.Sprintf("\\b%s\\b", substr), s)
	return match
}