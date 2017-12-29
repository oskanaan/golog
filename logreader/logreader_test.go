package logreader

import (
	"testing"
	"reflect"
)

func TestLogReader_parseLine(t *testing.T) {
	actual := "Test~Log~entry"
	expected := []string{"Test", "Log", "entry"}
    //Test parseLine directly
	result := parseLine(actual, "~")

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

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	result := *logReader.Tail()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Head(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"11/11/2010", "Thread-1", "com.test"},
		{"12/11/2010", "Thread-2", "com.test"},
		{"13/11/2010", "Thread-3", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	result := *logReader.Head()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_PageUp(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"13/11/2010", "Thread-3", "com.test"},
		{"14/11/2010", "Thread-4", "com.test"},
		{"15/11/2010", "Thread-5", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	result := *logReader.PageUp()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_PageUp_untilBeginning(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"11/11/2010", "Thread-1", "com.test"},
		{"12/11/2010", "Thread-2", "com.test"},
		{"13/11/2010", "Thread-3", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageUp()
	logReader.PageUp()
	logReader.PageUp()
	logReader.PageUp()
	logReader.PageUp()
	logReader.PageUp()
	result := *logReader.PageUp()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_PageDown(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"14/11/2010", "Thread-4", "com.test"},
		{"15/11/2010", "Thread-5", "com.test"},
		{"16/11/2010", "Thread-6", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageUp()
	logReader.PageUp()
	result := *logReader.PageDown()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_PageDown_untilEnd(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"16/11/2010", "Thread-6", "com.test"},
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageDown()
	logReader.PageDown()
	logReader.PageDown()
	logReader.PageDown()
	logReader.PageDown()
	logReader.PageDown()
	result := *logReader.PageDown()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Down(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"12/11/2010", "Thread-2", "com.test"},
		{"13/11/2010", "Thread-3", "com.test"},
		{"14/11/2010", "Thread-4", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageUp()
	logReader.PageUp()
	result := *logReader.Down()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Down_shouldntGoBeyondEndOfFile(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"16/11/2010", "Thread-6", "com.test"},
		{"17/11/2010", "Thread-7", "com.test"},
		{"18/11/2010", "Thread-8", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageUp()
	for i:=0 ; i<20 ; i++ {
		logReader.Down()
	}
	result := *logReader.Down()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Up(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"12/11/2010", "Thread-2", "com.test"},
		{"13/11/2010", "Thread-3", "com.test"},
		{"14/11/2010", "Thread-4", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageUp()
	result := *logReader.Up()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Up_shouldntGoBeforeBeginningOfFile(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"11/11/2010", "Thread-1", "com.test"},
		{"12/11/2010", "Thread-2", "com.test"},
		{"13/11/2010", "Thread-3", "com.test"},
	}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 3, ""})
	logReader.Tail()
	logReader.PageUp()
	for i:=0 ; i<20 ; i++ {
		logReader.Up()
	}
	result := *logReader.Up()

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

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}, 2, ""})
	result := *logReader.Tail()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Headers(t *testing.T) {
	input := "../test_logs/TestLogReader_Headers_input.log"
	expected := [] string {"Date", "Thread", "Package"}

	logReader := NewLogReader(input, Config{`~`, []string{"Date", "Thread", "Package"}, []int{15, 20, 10}, 3, ""})
	if !reflect.DeepEqual(logReader.GetHeaders(), expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, logReader.GetHeaders())
	}
}

func TestLogReader_GetColumnSizes(t *testing.T) {
	expected := []int{15, 20, 10}
    logReader := NewLogReader("", Config{`~`, []string{"Date", "Thread", "Package"}, expected, 3, ""})

    if !reflect.DeepEqual(logReader.GetColumnSizes(), expected) {
    	t.Errorf(`Expected column-sizes config to match the value returned by GetColumnSizes, expected %s, got %s`, expected, logReader.GetColumnSizes())
	}
}

func TestLogReader_GetSeverityColumnName(t *testing.T) {
	expected := "Test"
	logReader := NewLogReader("", Config{`~`, []string{"Date", "Thread", "Package", "Test"}, []int{1,2,3,4}, 3, "Test"})

	if logReader.GetSeverityColumnName() != expected {
		t.Errorf(`Expected column-sizes config to match the value returned by GetColumnSizes, expected %s, got %s`, expected, logReader.GetSeverityColumnName())
	}
}

func TestLogReader_lineCounter(t *testing.T){
	lines := "abc\n123\nefg\nmorelines\nend"
	expected := 5
	actual := countLines(lines)
	if expected != actual {
		t.Errorf("Expected %d lines, got %d", expected, actual)
	}
}