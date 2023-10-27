package main

import (
	"fmt"
	"strings"
)

func main() {
	var str1 string
	var str2 string
	fmt.Scan(&str1)
	fmt.Scan(&str2)

	charSlice1 := strings.Split(str1, "")
	charSlice2 := strings.Split(str2, "")

	charMap := map[string]int{}
	for _, v := range charSlice1 {
		charMap[v] = 0
		//charMap = append(charMap, v)
	}

	for _, v := range charSlice2 {
		_, exisCharMap := charMap[v]
		if exisCharMap {
			count := charMap[v]
			count++
			charMap[v] = count
		}
	}

	sum := 0
	for _, v := range charMap {
		sum = sum + v
	}

	fmt.Println(sum)
}
