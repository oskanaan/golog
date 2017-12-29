package logdisplay

import (
	"testing"
	"fmt"
)

func TestColorize_colorizeHeader(t *testing.T) {
	header := "header1\theader2\theader3"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2,4 ), header, "\033[0m")
	actual := colorizeHeader(header)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_errorLogEntryColoring(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 1, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, "ERROR")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_warnLogEntryColoring(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 3, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, "WARN")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_debugLogEntryColoring(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, "DEBUG")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_infoLogEntryColoring(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 6, 5), entry, "\033[0m")
	actual := colorizeLogEntry(entry, "INFO")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_traceLogEntryColoring(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 0, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, "TRACE")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_unknownSeverityDefaultsToTrace(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 0, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, "WEEEEWEEE!!")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}