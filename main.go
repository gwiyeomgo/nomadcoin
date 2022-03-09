package main

import (
	"github.com/gwiyeomgo/nomadcoin/rest"
)

func main() {
	//port 2개 실행시 go 사용
	rest.Start(3000)
	//go explorer.Start(4000)
}
