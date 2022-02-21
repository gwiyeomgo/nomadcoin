package main

import (
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/blockchain"
)

func main() {
	/*	chain := blockchain{}
		chain.addBlock("B block")
		chain.addBlock("S block")
		chain.addBlock("T block")
		chain.listBlocks()
	*/
	chain := blockchain.GetBlockchain()
	chain.AddBlock("Bcd")
	chain.AddBlock("Cde")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("Hash:%s\n", block.Hash)
		fmt.Printf("PreHash:%s\n", block.PreHash)
	}
}
