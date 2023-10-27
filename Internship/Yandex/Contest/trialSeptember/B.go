package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main20() {
	readNumbers := readFile("input.txt")
	sum := readNumbers[0] + readNumbers[1]
	//strOut := strconv.Itoa(sum)
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(strconv.Itoa(sum))
	if err2 != nil {
		log.Fatal(err2)
	}
}

func readFile(fileName string) []int {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var firstLine []int
	scanner := bufio.NewScanner(file)
	for i := 1; i < 2 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbersB(scanner.Text())
			//case 2:
			//secondLine = numbers(scanner.Text())
		}
	}

	return firstLine

}

func numbersB(s string) []int {
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
