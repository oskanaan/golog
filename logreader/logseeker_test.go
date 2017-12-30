package logreader

import (
	"os"
	"reflect"
	"testing"
)

func TestLogSeeker_readFileFromEnd_nonZeroOffset(t *testing.T) {
	input := "../test_logs/TestLogSeeker.log"
	expectedEntries := []string{"13/11/2010~Thread-3~com.test", "14/11/2010~Thread-4~com.test", "15/11/2010~Thread-5~com.test", "16/11/2010~Thread-6~com.test"}
	expectedOffset := 3
	file, _ := os.Open(input)
	actualEntries, actualOffset := readFileFromEnd(file, 4, 7)

	if !reflect.DeepEqual(expectedEntries, actualEntries) || expectedOffset != actualOffset {
		t.Errorf("Expected entries %s and offset %d, got %s and offset %d", expectedEntries, expectedOffset, actualEntries, actualOffset)
	}
}

func TestLogSeeker_readFileFromEnd_chainNonZeroOffsets(t *testing.T) {
	input := "../test_logs/TestLogSeeker.log"
	expectedEntries := []string{"13/11/2010~Thread-3~com.test", "14/11/2010~Thread-4~com.test"}
	expectedOffset := 3
	file, _ := os.Open(input)
	_, offset1 := readFileFromEnd(file, 2, 7)
	actualEntries, actualOffset := readFileFromEnd(file, 2, offset1)

	if !reflect.DeepEqual(expectedEntries, actualEntries) || expectedOffset != actualOffset {
		t.Errorf("Expected entries %s and offset %d, got %s and offset %d", expectedEntries, expectedOffset, actualEntries, actualOffset)
	}
}
