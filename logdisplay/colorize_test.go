package logdisplay

import (
	"fmt"
	"testing"
)

func TestColorize_colorizeHeader(t *testing.T) {
	header := "header1\theader2\theader3"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2, 4), header, "\033[0m")
	actual := colorizeHeader(header, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_errorLogEntryColoring(t *testing.T) {
	entry := "Blah\tERROR\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 1, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_errorLogEntryWithWarnTextColoring_shouldUseErrorColoring(t *testing.T) {
	entry := "Blah\tERROR\tBlah\tsome warn here"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 1, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_warnLogEntryColoring(t *testing.T) {
	entry := "Blah\tWARN\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 3, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_debugLogEntryColoring(t *testing.T) {
	entry := "Blah\tDEBUG\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 0, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_infoLogEntryColoring(t *testing.T) {
	entry := "Blah\tINFO\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 2, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_traceLogEntryColoring(t *testing.T) {
	entry := "Blah\tTRACE\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 6, 5), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestColorize_colorizeLogEntry_unknownSeverityDefaultsToTrace(t *testing.T) {
	entry := "Blah\tBlah\tBlah"
	expected := fmt.Sprint(fmt.Sprintf("\033[3%d;%d;1m", 0, 1), entry, "\033[0m")
	actual := colorizeLogEntry(entry, logdisplayConfig())

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
