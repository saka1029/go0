package main

func stringLength(b byte) int {
	switch {
	case b&0b10000000 == 0b00000000:
		return 1
	case b&0b11100000 == 0b11000000:
		return 2
	case b&0b11110000 == 0b11100000:
		return 3
	case b&0b11111000 == 0b11110000:
		return 4
	default:
		return -999
	}
}

type StringReader struct {
	s     string
	index int
}

func newStringReader(s string) *StringReader {
	return &StringReader{s, 0}
}

func (this *StringReader) Get() string {
	if this.index >= len(this.s) {
		return ""
	}
	length := stringLength(this.s[this.index])
	result := this.s[this.index : this.index+length]
	this.index += length
	return result
}
