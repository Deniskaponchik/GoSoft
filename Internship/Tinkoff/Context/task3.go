package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func task30() {
	var firstLine, secondLine, thirdline []int
	scanner := bufio.NewScanner(os.Stdin)
	for i := 1; i <= 3 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbers(scanner.Text())
		case 2:
			secondLine = numbers(scanner.Text())
		case 3:
			thirdline = numbers(scanner.Text())
		}
	}

	time := firstLine[1]
	floorTimeCol := thirdline[0]

	var result int
	if (floorTimeCol - secondLine[0]) < time {
		fmt.Println("if")
		fmt.Println(secondLine[firstLine[0]-1] - secondLine[0])
	} else {
		fmt.Println("else")
		fmt.Println(floorTimeCol - secondLine[0])
		fmt.Println(secondLine[firstLine[0]-1] - secondLine[0])
		result = (floorTimeCol - secondLine[0]) + (secondLine[firstLine[0]-1] - secondLine[0])
		fmt.Println(result)
	}

}

func numbers(s string) []int {
	var n []int
	for _, f := range strings.Fields(s) {
		i, err := strconv.Atoi(f)
		if err == nil {
			n = append(n, i)
		}
	}
	return n
}
