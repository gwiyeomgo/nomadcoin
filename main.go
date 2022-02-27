package main

import (
	"github.com/gwiyeomgo/nomadcoin/explorer"
	"github.com/gwiyeomgo/nomadcoin/rest"
)

func main() {
	go rest.Start(5000)
	explorer.Start(4000)
}
