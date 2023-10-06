package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var firstLine []int
	var secondLine []int
	scanner := bufio.NewScanner(os.Stdin)

	for i := 1; i < 3 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbersA(scanner.Text())
		case 2:
			secondLine = numbersA(scanner.Text())
		}
	}

	y1 := firstLine[0]
	y2 := secondLine[0]
	mo1 := firstLine[1]
	mo2 := secondLine[1]
	d1 := firstLine[2]
	d2 := secondLine[2]
	h1 := firstLine[3]
	h2 := secondLine[3]
	mi1 := firstLine[4]
	mi2 := secondLine[4]
	s1 := firstLine[5]
	s2 := secondLine[5]

	date1 := time.Date(y2, time.Month(mo2), d2, h2, mi2, s2, 0, time.UTC)
	date2 := time.Date(y1, time.Month(mo1), d1, h1, mi1, s1, 0, time.UTC)
	//date1 := time.Date(9009, 9, 11, 12, 21, 11, 0, time.UTC)
	//date2 := time.Date(1001, 5, 20, 14, 15, 16, 0, time.UTC)
	//date1 := time.Date(980, 3, 1, 10, 31, 37, 0, time.UTC)
	//date2 := time.Date(980, 2, 12, 10, 30, 1, 0, time.UTC)

	if date1.After(date2) {
		date1, date2 = date2, date1
	}
	days, hours, minutes, seconds := getDifference(date1, date2)
	fmt.Println(days, seconds+(minutes*60)+(hours*60*60))
}

func getDifference(a, b time.Time) (days, hours, minutes, seconds int) {
	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()
	totalDays1 := y1*365 + d1
	for i := 0; i < (int)(m1)-1; i++ {
		totalDays1 += monthDays[i]
	}
	totalDays2 := y2*365 + d2
	for i := 0; i < (int)(m2)-1; i++ {
		totalDays2 += monthDays[i]
	}
	days = totalDays2 - totalDays1
	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	return days, hours, minutes, seconds
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
func leapYears(date time.Time) (leaps int) {
	y, m, _ := date.Date()
	if m <= 2 {
		y--
	}
	leaps = y/4 + y/400 - y/100
	return leaps
}


func main11() {
	firstDate := time.Date(2022, 4, 13, 3, 0, 0, 0, time.UTC)
	secondDate := time.Date(2010, 2, 12, 6, 0, 0, 0, time.UTC)
	difference := firstDate.Sub(secondDate)
	fmt.Println(difference.Seconds())
}

func main10() {
	var firstLine []int
	var secondLine []int
	scanner := bufio.NewScanner(os.Stdin)

	for i := 1; i < 3 && scanner.Scan(); i++ {
		switch i {
		case 1:
			firstLine = numbersA(scanner.Text())
		case 2:
			secondLine = numbersA(scanner.Text())
		}
	}

	y1 := firstLine[0]
	y2 := secondLine[0]
	mo1 := firstLine[1]
	mo2 := secondLine[1]
	d1 := firstLine[2]
	d2 := secondLine[2]
	h1 := firstLine[3]
	h2 := secondLine[3]
	mi1 := firstLine[4]
	mi2 := secondLine[4]
	s1 := firstLine[5]
	s2 := secondLine[5]

	//(9009 - 1001) * 365 =
	//difYearDays := (secondLine[0] - firstLine[0]) * 365
	countDays := (y2 - y1) * 365
	//+
	//daysInMonths := {0, }
	daysInMonths := [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	//countDays := 0
	firstComing := 1
	j := true
	for i := mo2; j == true; i-- {
		if i == 0 {
			i = 13
		} else {
			if i == mo1 {
				if mo2 == mo1 {
					countDays = countDays + (d2 - d1)
					j = false
					//fmt.Println(1, countDays)
				} else {
					if h2 < h1 {
						countDays = countDays + daysInMonths[i] - d1 - 1
					} else {
						countDays = countDays + daysInMonths[i] - d1 //- 1
					}
					j = false
					//fmt.Println(2, countDays)
				}
			} else {
				// i != mo1
				if firstComing == 1 {
					countDays = countDays + d2
					firstComing = 0
					//fmt.Println(3, countDays)
				} else {
					countDays = countDays + daysInMonths[i]
					//fmt.Println(4, countDays)
				}
			}
		}
		//
	}
	//fmt.Println(countDays)

	//secInDay := 86400
	countHours := 0
	countMinutes := 0
	countSeconds := 0

	if h2 < h1 {
		countHours = 24 - (h1 - h2)
	} else {
		countHours = h2 - h1
	}
	//fmt.Println(countHours)

	if mi2 > mi1 {
		if s2 < s1 {
			countMinutes = mi2 - mi1 - 1
		} else {
			countMinutes = mi2 - mi1
		}
	} else {
		countMinutes = mi2 + (60 - mi1)
	}
	//fmt.Println(countMinutes)

	if s2 < s1 {
		countSeconds = s2 + (60 - s1)
	} else {
		countSeconds = s2 - s1
	}
	//fmt.Println(countSeconds)

	countSeconds = countSeconds + (countMinutes * 60) + (countHours * 60 * 60)

	//result := [2]int{countDays, countSeconds}
	//fmt.Println(result)
	fmt.Println(strconv.Itoa(countDays) + " " + strconv.Itoa(countSeconds))

}
*/
