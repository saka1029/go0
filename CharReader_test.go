package main

import (
	"testing"
)

func TestCharReader(t *testing.T) {
	m := "TestCharReader"
	assertEquals(t, m, "a", 1, CharLength("a"[0]))
	assertEquals(t, m, "ä", 2, CharLength("ä"[0]))
	assertEquals(t, m, "あ", 3, CharLength("あ"[0]))
	assertEquals(t, m, "𩸽", 4, CharLength("𩸽"[0]))
	cr := NewCharReader("aäあ𩸽")
	assertEquals(t, m, "a", "a", cr.Get())
	assertEquals(t, m, "ä", "ä", cr.Get())
	assertEquals(t, m, "あ", "あ", cr.Get())
	assertEquals(t, m, "𩸽", "𩸽", cr.Get())
	assertEquals(t, m, "", "", cr.Get())
}
