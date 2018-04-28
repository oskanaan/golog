package logreader

import (
	"reflect"
	"testing"
)

func logreaderConfig(file string, sizes []int) LogReaderConfig{
	return LogReaderConfig{
		Files: []LogFile{{file, "Name"}},
		Seperator: "~",
		Headers: []Header{
			{"Date", sizes[0]},
			{"Thread", sizes[1]},
			{"Package", sizes[2]},
		},
	}
}

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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
	logReader.Tail()
	logReader.PageUp()
	for i := 0; i < 20; i++ {
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
	logReader.Tail()
	logReader.PageUp()
	result := *logReader.Up()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Up_WithStackTrace(t *testing.T) {
	input := "../test_logs/TestLogReader_Navigate_withstacktrace.log"
	expected := [][]string{
		{"14/11/2010", "Thread-8", "com.test"},
	}

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
	logReader.Head()
	logReader.PageDown()
	logReader.PageDown()
	logReader.PageDown()
	logReader.PageDown()
	logReader.Down()
	logReader.Up()
	logReader.Up()
	logReader.Up()
	logReader.Up()
	logReader.Up()
	result := *logReader.Up()

	if !reflect.DeepEqual(result[0], expected[0]) {
		t.Errorf(`Output Log: Expected %s got %s`, expected[0], result[0])
	}
}

func TestLogReader_Up_shouldntGoBeforeBeginningOfFile(t *testing.T) {
	input := "../test_logs/TestLogReader_Tail_input.log"
	expected := [][]string{
		{"11/11/2010", "Thread-1", "com.test"},
		{"12/11/2010", "Thread-2", "com.test"},
		{"13/11/2010", "Thread-3", "com.test"},
	}

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
	logReader.Tail()
	logReader.PageUp()
	for i := 0; i < 20; i++ {
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

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(2)
	result := *logReader.Tail()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_Headers(t *testing.T) {
	input := "../test_logs/TestLogReader_Headers_input.log"
	expected := []string{"Date", "Thread", "Package"}

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(3)
	if !reflect.DeepEqual(logReader.GetHeaders(), expected) {
		t.Errorf(`Output Log: Expected %s got %s`, expected, logReader.GetHeaders())
	}
}

func TestLogReader_GetColumnSizes(t *testing.T) {
	expected := []int{15, 20, 10}
	logReader := NewLogReader(logreaderConfig("", expected))
	logReader.SetCapacity(3)

	if !reflect.DeepEqual(logReader.GetColumnSizes(), expected) {
		t.Errorf(`Expected column-sizes config to match the value returned by GetColumnSizes, expected %s, got %s`, expected, logReader.GetColumnSizes())
	}
}

func TestLogReader_lineCounter(t *testing.T) {
	lines := "abc\n123\nefg\nmorelines\nend"
	expected := 5
	actual := countLines(lines)
	if expected != actual {
		t.Errorf("Expected %d lines, got %d", expected, actual)
	}
}

func TestLogReader_Search_ShouldReturnFirstInstance(t *testing.T) {
	input := "../test_logs/TestLogReader_Search.log"
	expected := [][]string{{"       at rx.Observable$31.onError(Observable.java:7204)"}, {"       at rx.observers.SafeSubscriber._onError(SafeSubscriber.java:127)"}}

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(2)
	actual, _ := logReader.Search("(Observable.java:7204)", 0)

	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expected %s lines, got %s", expected, *actual)
	}
}

func TestLogReader_Search_NextShouldReturnNextInstance(t *testing.T) {
	input := "../test_logs/TestLogReader_Search.log"
	expected := [][]string{{"Caused by: rx.exceptions.MissingBackpressureException"}, {"       at rx.internal.util.RxRingBuffer.onNext(RxRingBuffer.java:222)"}}

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(2)
	logReader.Search("Caused by:", 0)
	actual, _ := logReader.Search("Caused by:", 0)

	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expected %s lines, got %s", expected, *actual)
	}
}

func TestLogReader_Search_NextAtEndShouldReturnToBeginning(t *testing.T) {
	input := "../test_logs/TestLogReader_Search.log"
	expected := [][]string{{"Caused by: rx.exceptions.OnErrorNotImplementedException"}, {"       at rx.Observable$31.onError(Observable.java:7204)"}}

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(2)
	logReader.Search("Caused by:", 0)
	logReader.Search("Caused by:", 0)
	actual, _ := logReader.Search("Caused by:", 0)

	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("Expected %s lines, got %s", expected, *actual)
	}
}

func TestLogReader_Progress_Tail(t *testing.T) {
	input := "../test_logs/TestLogReader_Progress.log"

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(2)
	logReader.Tail()
	actual := logReader.Progress()
	expected := 100

	if actual != expected {
		t.Errorf("Expected %s lines, got %s", expected, actual)
	}
}

func TestLogReader_Progress_Head(t *testing.T) {
	input := "../test_logs/TestLogReader_Progress.log"

	logReader := NewLogReader(logreaderConfig(input, []int{10, 10, 10}))
	logReader.SetCapacity(2)
	logReader.Head()
	actual := logReader.Progress()
	expected := 58*100/291

	if actual != expected {
		t.Errorf("Expected %s lines, got %s", expected, actual)
	}
}