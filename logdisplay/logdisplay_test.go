package logdisplay

import (
	"bytes"
	"github.com/oskanaan/golog/logreader"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"text/tabwriter"
)

func TestLogReader_Tail(t *testing.T) {
	expected := [][]string{
		{"16/11/2010", "Thread-6", "com.test"},
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := logreader.NewLogReader(logreaderConfig("../test_logs/TestLogDisplay_Tail_input.log", []int{10, 10, 10}))
	logReader.SetCapacity(3)
	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	logdisplay.tail()
	result := *logdisplay.currentPage

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Tail_3LinesLog_WithCapacitySizeEquals2(t *testing.T) {
	expected := [][]string{
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := logreader.NewLogReader(logreaderConfig("../test_logs/TestLogDisplay_Tail_input.log", []int{10, 10, 10}))
	logReader.SetCapacity(2)
	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	logdisplay.tail()
	result := logdisplay.currentPage

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, *result)
	}
}

func TestLogDisplay_formatColumnText_positiveColumnSize(t *testing.T) {
	logReader := logreader.NewLogReader(logreaderConfig("", []int{10, 10, 10}))
	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	expected := "characters"
	actual := logdisplay.formatColumnText("More than 10 characters", 1)

	if expected != actual {
		t.Errorf(`Output Log: Expected %s got %s`, expected, actual)
	}
}

func TestLogDisplay_formatColumnText_negativeColumnSize(t *testing.T) {
	logReader := logreader.NewLogReader(logreaderConfig("", []int{-1, 10, 10}))
	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	expected := "More than 10 characters"
	actual := logdisplay.formatColumnText("More than 10 characters", 0)

	if expected != actual {
		t.Errorf(`Output Log: Expected %s got %s`, expected, actual)
	}
}

func TestLogDisplay_formatColumnText_zeroColumnSize(t *testing.T) {
	logReader := logreader.NewLogReader(logreaderConfig("", []int{0, 10, 10}))
	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	expected := ""
	actual := logdisplay.formatColumnText("More than 10 characters", 0)

	if expected != actual {
		t.Errorf(`Output Log: Expected %s got %s`, expected, actual)
	}
}

func TestLogDisplay_formatColumnText_testWhenColumnSizeIsLessThanConfiguredValue(t *testing.T) {
	logReader := logreader.NewLogReader(logreaderConfig("", []int{15, 10, 10}))
	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	expected := "7 chars        "
	actual := logdisplay.formatColumnText("7 chars", 0)

	if expected != actual {
		t.Errorf(`Output Log: Expected %s got %s`, expected, actual)
	}
}

func TestLogDisplay_writeHeader(t *testing.T) {
	logReader := logreader.NewLogReader(logreaderConfig("", []int{10, 10, 10}))
	var actual bytes.Buffer
	expectedRegexp := "ate.*\\wd\\w.*Package\n"

	logdisplay := NewLogDisplay(&logReader, logdisplayConfig())
	tabWriter := new(tabwriter.Writer)
	tabWriter.Init(&actual, 0, 8, 0, '\t', tabwriter.AlignRight)
	logdisplay.writeHeader(tabWriter)
	tabWriter.Flush()

	if match, _ := regexp.MatchString(expectedRegexp, strings.TrimSpace(actual.String())); match {
		t.Errorf(`Output Log: Expected %s to match %s`, actual.String(), expectedRegexp)
	}
}
