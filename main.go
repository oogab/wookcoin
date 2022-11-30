package main

import "github.com/nomadcoders/nomadcoin/blockchain"

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("Second Block")
}
