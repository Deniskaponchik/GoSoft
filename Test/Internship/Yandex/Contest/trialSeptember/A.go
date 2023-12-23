package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main10() {
	var firstLine []int
	scanner := bufio.NewScanner(os.Stdin)

	for i := 1; i < 2 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbersA(scanner.Text())
			//case 2:
			//secondLine = numbers(scanner.Text())
		}
	}

	fmt.Println(firstLine[0] + firstLine[1])

}

func numbersA(s string) []int {
	var n []int
	for _, f := range strings.Fields(s) {
		i, err := strconv.Atoi(f)
		if err == nil {
			n = append(n, i)
		}
	}
	return n
}

/*
Даны два числа A и B. Вам нужно вычислить их сумму A+B.
В этой задаче для работы с входными и выходными данными вы можете использовать и файлы и потоки на ваше усмотрение.
*/
