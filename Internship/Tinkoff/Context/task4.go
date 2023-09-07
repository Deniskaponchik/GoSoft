package main

import (
	"bufio"
	"os"
)

func task40() {
	var firstLine, secondLine []int
	scanner := bufio.NewScanner(os.Stdin)
	for i := 1; i <= 3 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbers(scanner.Text())
		case 2:
			secondLine = numbers(scanner.Text())
		}
	}

	var sum1 int
	for i := 0; i < firstLine[0]; i++ {
		sum1 = secondLine[i] + sum1
	}

}
