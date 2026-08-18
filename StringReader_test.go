package main

import (
	"testing"
)

func TestStringReader(t *testing.T) {
	m := "TestStringReader"
	assertEquals(t, m, "a", 1, stringLength("a"[0]))
	assertEquals(t, m, "ä", 2, stringLength("ä"[0]))
	assertEquals(t, m, "あ", 3, stringLength("あ"[0]))
	assertEquals(t, m, "𩸽", 4, stringLength("𩸽"[0]))
	sr := NewStringReader("aäあ𩸽")
	assertEquals(t, m, "a", "a", sr.Get())
	assertEquals(t, m, "ä", "ä", sr.Get())
	assertEquals(t, m, "あ", "あ", sr.Get())
	assertEquals(t, m, "𩸽", "𩸽", sr.Get())
	assertEquals(t, m, "", "", sr.Get())
}
