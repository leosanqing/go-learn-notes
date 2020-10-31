package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"runtime"
	"time"
)

func euler() {
	fmt.Println(cmplx.Pow(math.E, 1i*math.Pi) + 1)
}
func main() {



	var a [10]int
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				a[i]++
			}
		}(i)
	}
	time.Sleep(time.Millisecond)
	fmt.Println(a)


	mapDemo()
	extendSlice()
	fmt.Println(runtime.GOARCH)
	euler()

	array := [...]int{0, 1, 2, 3, 4, 5, 6}
	fmt.Println(array[1:2])
}

func extendSlice() {
	array := [...]int{0, 1, 2, 3, 4, 5, 6, 7}

	s1 := array[2:6]
	s2 := s1[3:5]

	fmt.Printf("s1 = %v, len(s1) = %d ,cap(s1) = %d\n", s1, len(s1), cap(s1))
	fmt.Println(s2)
}

func mapDemo() {
	m1 := map[string]string{
		"name":   "张三",
		"age":    "18",
		"school": "swu",
	}

	fmt.Print(m1)

	name := m1["name"]
	fmt.Println("name:" + name)

	if nema, ok := m1["nema"]; ok {
		fmt.Println("nema:=", nema)
	} else {
		fmt.Println("key does not exist")
	}

}
