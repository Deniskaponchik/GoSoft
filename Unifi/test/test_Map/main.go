package main

import "fmt"

type Ap struct {
	Count int
}

func new(count int) *Ap {
	return &Ap{
		Count: count,
	}
}
func one() map[string]*Ap {
	ap1 := new(1)
	ap2 := new(2)
	map1 := make(map[string]*Ap)
	map1["mac1"] = ap1
	map1["mac2"] = ap2
	return map1
}
func two(map1 map[string]*Ap) {
	k, exis := map1["mac1"]
	if exis {
		k.Count = 3
	}
}

func main() {
	map1 := one()

	two(map1)
	for k, v := range map1 {
		fmt.Println(k, v.Count)
	}

}
