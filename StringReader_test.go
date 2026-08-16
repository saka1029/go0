package main

import (
	"fmt"
	"testing"
)

func stringReaderTest(t *testing.T, s string, expected int) {
	actual := stringLength(s[0])
	fmt.Printf("%s length=%d\n", s, actual)
	if expected != actual {
		t.Errorf("TestString %s -> %d not %d", s, actual, expected)
	}
}
func stringReaderTest2(t *testing.T, sr *StringReader, expected string) {
	actual := sr.Get()
	fmt.Printf("%s expected=%s\n", actual, expected)
	if expected != actual {
		t.Errorf("TestString %s not %s", actual, expected)
	}
}

func TestStringReader(t *testing.T) {
	stringReaderTest(t, "a", 1)
	stringReaderTest(t, "ä", 2)
	stringReaderTest(t, "あ", 3)
	stringReaderTest(t, "𩸽", 4)
	sr := NewStringReader("aäあ𩸽")
	stringReaderTest2(t, sr, "a")
	stringReaderTest2(t, sr, "ä")
	stringReaderTest2(t, sr, "あ")
	stringReaderTest2(t, sr, "𩸽")
	stringReaderTest2(t, sr, "")
}
