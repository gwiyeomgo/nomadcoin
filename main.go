package main

import (
	"fmt"
)

func plus(array ...int) int {
	// a는 int의 array
	results := 0
	for _, item := range array {
		results += item
	}
	return results
}
func main() {
	result := plus(1, 1, 1, 1, 1, 2)
	fmt.Println(result)
	name := "gwiyeom go~~!!#!#!"
	for _, item := range name {
		//fmt.Println(item) //byte 타입으로 출력된다.
		fmt.Println(string(item))
	}
}
