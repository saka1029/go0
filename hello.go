package main

import (
	"fmt"
)

func main() {
	var x = [][]int{{1, 2}, {3, 4, 5}}
	// x[0][0] = 10
	// x[1] = []int{40, 50}
	var y = [][]int{{1, 2}, {3, 4, 5}}
	fmt.Println(x, y)
}
