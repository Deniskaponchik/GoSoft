package main

import "fmt"

func main() {

	sliceInt100 := [100]int{}
	for i := 0; i < 100; i++ {
		sliceInt100[i] = i
	}

	//sliceInt := [len(sliceInt100)]int{}
	sliceInt := make([]int, int(len(sliceInt100)))
	j := 0
	for i := len(sliceInt100) - 1; i > -1; i-- {
		sliceInt[j] = sliceInt100[i]
		j++
	}

	lenSliceInt := len(sliceInt) - 1
	for i := 0; i < lenSliceInt; i++ {
		fmt.Print(sliceInt[i])
	}
}
