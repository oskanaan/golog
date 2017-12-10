package logdisplay

import (
	"testing"
	"github.com/oskanaan/golog/logreader"
	"reflect"
)

func TestLogReader_Tail(t *testing.T) {
	input := "../test_logs/TestLogDisplay_Tail_input.log"
	expected := [][]string{
		{"16/11/2010", "Thread-6", "com.test"},
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := logreader.NewLogReader(input, logreader.Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3})
	logdisplay := NewLogDisplay(logReader)
	result := logdisplay.Tail()

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Tail_3LinesLog_WithCapacitySizeEquals2(t *testing.T) {
	input := "../test_logs/TestLogDisplay_Tail_input.log"
	expected := [][]string{
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := logreader.NewLogReader(input, logreader.Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 2})
	logdisplay := NewLogDisplay(logReader)
	result := logdisplay.Tail()

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}