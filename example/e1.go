package main

import (
	"fmt"

	"github.com/starjun/gotools"
)

func main() {
	a := []int{1, 2, 1, 3, 4}
	b := []int{4, 3, 5, 6}

	c := gotools.GetDifferentIntArray(a, b)

	fmt.Println(c)

}
