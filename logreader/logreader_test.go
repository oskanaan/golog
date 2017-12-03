package logreader

import (
	"bufio"
	"bytes"
	"testing"
	"time"
	"io/ioutil"
	"strings"
)

func TestLogReader_parseLine(t *testing.T) {
	var output bytes.Buffer
	actual := "Test~Log~entry"
	expected := "Test      Log       entry     \n"

	writer := bufio.NewWriter(&output)
	logReader := NewLogReader(actual, *writer, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}})
	//Test parseLine directly
	isReadSuccessful := logReader.parseLine(actual)

	if !isReadSuccessful {
		t.Errorf(`Output Log: Expected a line to be read, bug got an error`)
	}

	result := output.String()
	if result != expected {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}

func TestLogReader_StartReading(t *testing.T) {
	var output bytes.Buffer
	input := "../test_logs/TestLogReader_StartReading_input.log"
	expectedBytes,_ := ioutil.ReadFile("../test_logs/TestLogReader_StartReading_output.log")

	writer := bufio.NewWriter(&output)
	logReader := NewLogReader(input, *writer, Config{`~`, []string{"Date", "Thread", "Package"}, []int{10, 10, 10}})
	go logReader.StartReading()
	time.Sleep(1 * time.Second)

	result := output.String()
	if result != string(expectedBytes) {
		t.Errorf(`Output Log: Expected %s got %s`, string(expectedBytes), result)
	}
}

func TestLogReader_Headers(t *testing.T) {
	var output bytes.Buffer
	input := "../test_logs/TestLogReader_Headers_input.log"
	expectedBytes,_ := ioutil.ReadFile("../test_logs/TestLogReader_Headers_output.log")

	writer := bufio.NewWriter(&output)
	logReader := NewLogReader(input, *writer, Config{`~`, []string{"Date", "Thread", "Package"}, []int{15, 20, 10}})
	go logReader.StartReading()
	time.Sleep(1 * time.Second)

	result := output.String()
	result = result[:strings.Index(result, "\n")]
	expected := string(expectedBytes)[:strings.Index(string(expectedBytes), "\n")]
	if result != expected {
		t.Errorf(`Output Log: Expected %s got %s`, expected, result)
	}
}
