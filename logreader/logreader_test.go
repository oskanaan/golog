package logreader

import (
	"testing"
	"reflect"
)

func TestLogReader_parseLine(t *testing.T) {
	actual := "Test~Log~entry"
	expected := []string{"Test", "Log", "entry"}

	logReader := NewLogReader(actual, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3})
	//Test parseLine directly
	result := logReader.parseLine(actual)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Tail(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"16/11/2010", "Thread-6", "com.test"},
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3})
	result := logReader.Tail()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Tail_3LinesLog_WithCapacitySizeEquals2(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 2})
	result := logReader.Tail()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Headers(t *testing.T) {
	input := "../test_logs/TestLogReader_Headers_input.log"
	expected := [] string {"Date", "Thread", "Package"}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{15, 20, 10}, 3})
	if !reflect.DeepEqual(logReader.GetHeaders(), expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, logReader.GetHeaders())
	}
}

func TestLogReader_GetColumnSizes(t *testing.T) {
	expected := []int{15, 20, 10}
    logReader := NewLogReader("", Config{`~`, []string{"Date", "Thread", "Package"}, expected, 3})

    if !reflect.DeepEqual(logReader.GetColumnSizes(), expected) {
    	t.Errorf(`Expected column-sizes config to match the value returned by GetColumnSizes, expected %s, got %s`, expected, logReader.GetColumnSizes())
	}
}

func TestLogReader_ScrollTo(t *testing.T){

}