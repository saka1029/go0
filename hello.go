package main

import (
	"fmt"
)

func main() {
	s := []string{"first", "second", "third"}
	fmt.Println(s, len(s)) // [first second third] 3
	clear(s)
	fmt.Println(s, len(s))                                            // [ ] 3
	fmt.Printf("s[0]=|%s|, s[1]=|%s|, s[2]=|%s|\n", s[0], s[1], s[2]) // s[0]=||, s[1]=||, s[2]=||

}
