package logdisplay

import (
	"fmt"
	"regexp"
)

//prepends/appends the header color to the log headers
//Returns a color coded string
func colorizeHeader(header string, logdisplayConfig LogDisplayConfig) string {
	return fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2, 4), header, "\033[0m")
}

//prepends/appends the color of the log entry based on the severity
//Returns a color coded string
//Todo: make the colors and severity values configurable
func colorizeLogEntry(logEntry string, logdisplayConfig LogDisplayConfig) string {
	colorCode := severityColorCode(logEntry, logdisplayConfig)
	pre := fmt.Sprintf("\033[3%d;%d;1m", colorCode...)
	return fmt.Sprint(pre, logEntry, "\033[0m")
}

func severityColorCode(entry string, logdisplayConfig LogDisplayConfig) []interface{} {
	for _, code := range logdisplayConfig.Severities {
		if severityMatch(entry, code.Severity) {
			return code.Colors
		}
	}

	return logdisplayConfig.Severities[len(logdisplayConfig.Severities) - 1].Colors
}

func severityMatch(logEntry, severityRegex string) bool {
	match, _ := regexp.MatchString(severityRegex, logEntry)
	return match
}