package main

import (
	"fmt"

	"github.com/oogab/wookchain/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain()
	fmt.Println(chain)
}
