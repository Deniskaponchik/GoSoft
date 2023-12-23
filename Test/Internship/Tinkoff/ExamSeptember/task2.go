package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string
	fmt.Scan(&str)
	res := strings.Split(str, "")

	x2slice := make([][7]int, 1)
	sheriff := [7]int{}

	countS := 0
	countH := 0
	countE := 0
	countR := 0
	countI := 0
	countF := 0
	countF2 := 0
	var lenX2slice int

	for _, v := range res {
		lenX2slice = len(x2slice)

		if v == "s" {
			if countS != lenX2slice {
				x2slice[countS][0] = 1
			} else {
				sheriff = [7]int{1, 0, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countS++
		} else if v == "h" {
			if countH != lenX2slice {
				x2slice[countH][1] = 1
			} else {
				sheriff = [7]int{0, 1, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countH++
		} else if v == "e" {
			if countE != lenX2slice {
				x2slice[countE][2] = 1
			} else {
				sheriff = [7]int{0, 0, 1, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countE++
		} else if v == "r" {
			if countR != lenX2slice {
				x2slice[countR][3] = 1
			} else {
				sheriff = [7]int{0, 0, 0, 1, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countR++
		} else if v == "i" {
			if countI != lenX2slice {
				x2slice[countI][4] = 1
			} else {
				sheriff = [7]int{0, 0, 0, 0, 1, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countI++
		} else if v == "f" {
			countF2 = countF % 2
			if countF2/2 != lenX2slice {
				if countF2 == 0 {
					x2slice[countF/2][5] = 1
				} else {
					x2slice[countF/2][6] = 1
				}
			} else {
				if countF2 != 0 {
					sheriff = [7]int{0, 0, 0, 0, 0, 1, 0}
				} else {
					sheriff = [7]int{0, 0, 0, 0, 0, 0, 1}
				}
				x2slice = append(x2slice, sheriff)
			}
			countF++
		}
	}

	countSheriff := 0
	var sum int
	for i := 0; i < lenX2slice; i++ {
		sum = 0
		for _, v := range x2slice[i] {
			if v != 0 {
				sum++
			}
		}
		if sum == 7 {
			countSheriff++
		}
	}
	fmt.Println(countSheriff)

}

/*
func main22() {
	var str string
	fmt.Scan(&str)
	res := strings.Split(str, "")
	//fmt.Printf("Data type: %T\n", res[0])
	//fmt.Printf("Characters: %q\n", res)

	//x2slice := make([][]string , 1)
	x2slice := make([][7]int, 1)
	//sheriff := []string{}
	sheriff := [7]int{}

	countS := 0
	countH := 0
	countE := 0
	countR := 0
	countI := 0
	countF := 0
	countF2 := 0
	var lenX2slice int

	for _, v := range res {
		lenX2slice = len(x2slice)
		fmt.Println("len(x2slice) = " + strconv.Itoa(int(lenX2slice)))
		if v == "s" {
			if countS != lenX2slice {
				x2slice[countS][0] = 1
				fmt.Println("s if")
			} else {
				sheriff = [7]int{1, 0, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
				fmt.Println("s else")
			}
			countS++
		} else if v == "h" {
			if countH != lenX2slice {
				x2slice[countH][1] = 1
				fmt.Println("h if")
			} else {
				sheriff = [7]int{0, 1, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
				fmt.Println("h else")
			}
			countH++
		} else if v == "e" {
			if countE != lenX2slice {
				x2slice[countE][2] = 1
				fmt.Println("e if")
			} else {
				sheriff = [7]int{0, 0, 1, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
				fmt.Println("e else")
			}
			countE++
		} else if v == "r" {
			if countR != lenX2slice {
				x2slice[countR][3] = 1
				fmt.Println("r if")
			} else {
				sheriff = [7]int{0, 0, 0, 1, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
				fmt.Println("r else")
			}
			countR++
		} else if v == "i" {
			if countI != lenX2slice {
				x2slice[countI][4] = 1
				fmt.Println("i if")
			} else {
				sheriff = [7]int{0, 0, 0, 0, 1, 0, 0}
				x2slice = append(x2slice, sheriff)
				fmt.Println("i else")
			}
			countI++
		} else if v == "f" {
			countF2 = countF % 2
			if countF2/2 != lenX2slice {
				if countF2 == 0 {
					x2slice[countF/2][5] = 1
					fmt.Println("f if if")
				} else {
					x2slice[countF/2][6] = 1
					fmt.Println("f if else")
				}
			} else {
				if countF2 != 0 {
					sheriff = [7]int{0, 0, 0, 0, 0, 1, 0}
					fmt.Println("f else if")
				} else {
					sheriff = [7]int{0, 0, 0, 0, 0, 0, 1}
					fmt.Println("f if else")
				}
				x2slice = append(x2slice, sheriff)
			}
			countF++
		}
	}
	fmt.Println(lenX2slice)
	fmt.Println(x2slice)

	countSheriff := 0
	var sum int
	for i := 0; i < lenX2slice; i++ {
		sum = 0
		for _, v := range x2slice[i] {
			if v != 0 {
				sum++
				fmt.Println("sum++")
			}
		}
		if sum == 7 {
			countSheriff++
		}
	}
	fmt.Println(countSheriff)

}

func main21() {
	var str string
	fmt.Scan(&str)
	res := strings.Split(str, "")

	x2slice := make([][7]int, 1)
	sheriff := [7]int{}

	countS := 0
	countH := 0
	countE := 0
	countR := 0
	countI := 0
	countF := 0
	countF2 := 0
	var lenX2slice int

	for _, v := range res {
		lenX2slice = len(x2slice)

		if v == "s" {
			if countS != lenX2slice {
				x2slice[countS][0] = 1
			} else {
				sheriff = [7]int{1, 0, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countS++
		} else if v == "h" {
			if countH != lenX2slice {
				x2slice[countH][1] = 1
			} else {
				sheriff = [7]int{0, 1, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countH++
		} else if v == "e" {
			if countE != lenX2slice {
				x2slice[countE][2] = 1
			} else {
				sheriff = [7]int{0, 0, 1, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countE++
		} else if v == "r" {
			if countR != lenX2slice {
				x2slice[countR][3] = 1
			} else {
				sheriff = [7]int{0, 0, 0, 1, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countR++
		} else if v == "i" {
			if countI != lenX2slice {
				x2slice[countI][4] = 1
			} else {
				sheriff = [7]int{0, 0, 0, 0, 1, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countI++
		} else if v == "f" {
			countF2 = countF % 2
			if countF2/2 != lenX2slice {
				if countF2 == 0 {
					x2slice[countF/2][5] = 1
				} else {
					x2slice[countF/2][6] = 1
				}
			} else {
				if countF2 != 0 {
					sheriff = [7]int{0, 0, 0, 0, 0, 1, 0}
				} else {
					sheriff = [7]int{0, 0, 0, 0, 0, 0, 1}
				}
				x2slice = append(x2slice, sheriff)
			}
			countF++
		}
	}

	countSheriff := 0
	var sum int
	for i := 0; i < lenX2slice; i++ {
		sum = 0
		for _, v := range x2slice[i] {
			if v != 0 {
				sum++
			}
		}
		if sum == 7 {
			countSheriff++
		}
	}
	fmt.Println(countSheriff)

}

func main20() {
	var str string
	fmt.Scan(&str)
	res := strings.Split(str, "")
	//fmt.Printf("Data type: %T\n", res[0])
	//fmt.Printf("Characters: %q\n", res)

	//x2slice := make([][]string , 1)
	x2slice := make([][7]int, 1)
	//sheriff := []string{}
	sheriff := [7]int{}

	countS := 0
	countH := 0
	countE := 0
	countR := 0
	countI := 0
	countF := 0
	var lenX2slice int

	for _, v := range res {
		lenX2slice = len(x2slice)
		if v == "s" {
			if countS != lenX2slice {
				x2slice[countS][0] = 1
			} else {
				sheriff = [7]int{1, 0, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
			countS++
		} else if v == "h" {
			if countH == lenX2slice {
				x2slice[lenX2slice][1] = 1
				countH++
			} else {
				sheriff = [7]int{0, 1, 0, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
		} else if v == "e" {
			if countE == lenX2slice {
				x2slice[lenX2slice][2] = 1
				countE++
			} else {
				sheriff = [7]int{0, 0, 1, 0, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
		} else if v == "r" {
			if countR == lenX2slice {
				x2slice[lenX2slice][3] = 1
				countR++
			} else {
				sheriff = [7]int{0, 0, 0, 1, 0, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
		} else if v == "i" {
			if countI == lenX2slice {
				x2slice[lenX2slice][4] = 1
				countI++
			} else {
				sheriff = [7]int{0, 0, 0, 0, 1, 0, 0}
				x2slice = append(x2slice, sheriff)
			}
		} else if v == "f" {
			if countF == lenX2slice {
				if countF%2 != 0 {
					x2slice[lenX2slice][5] = 1
				} else {
					x2slice[lenX2slice][6] = 1
				}
				countF++
			} else {
				if countF%2 != 0 {
					sheriff = [7]int{0, 0, 0, 0, 0, 1, 0}
				} else {
					sheriff = [7]int{0, 0, 0, 0, 0, 0, 1}
				}
				x2slice = append(x2slice, sheriff)
			}
		}
	}
	fmt.Println(lenX2slice)

	countSheriff := 0
	sum := 0
	for i := 0; i < lenX2slice; i++ {
		for _, v := range x2slice[i] {
			if v != 0 {
				sum++
			}
		}
		if sum == 7 {
			countSheriff++
		}
	}
	fmt.Println(countSheriff)

}
*/
