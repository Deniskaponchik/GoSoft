package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main1() {
	var firstLine, secondLine []int
	scanner := bufio.NewScanner(os.Stdin)
	for i := 1; i <= 2 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbers(scanner.Text())
		case 2:
			secondLine = numbers(scanner.Text())
		}
	}

	fmt.Println(findPurchase(secondLine, firstLine[1]))

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

func findPurchase(revs []int, dollars int) (purchase int) {
	count := 0
	for _, v := range revs {
		if dollars > v && purchase < v {
			purchase = v
			count++
		}
	}
	if count == 0 {
		purchase = 0
	}
	return
}

/*
// func findMinAndMax(a [5]int) (min int, max int) {
func findMinAndMax(revs []int) (min int, max int) {
	min = revs[0]
	max = revs[0]
	for _, value := range revs {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func var2() {
	var rev, dol int
	if _, err := fmt.Scan(&rev, &dol); err != nil {
		log.Print(" Scan for i, j & k failed, due to ", err)
		return
	}

	//fmt.Println(i)
	//fmt.Println(j)
}

func var1() {
	var rev int
	var dol int

	fmt.Scanf("%d %e", &rev, &dol)

	//fmt.Scan(&a)
	//fmt.Scan(&b)

	fmt.Println(rev)
	fmt.Println(dol)
}
*/
