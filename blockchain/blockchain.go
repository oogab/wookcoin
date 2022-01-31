package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block
}

var b *blockchain
var once sync.Once

func (b *block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}

// func (b *blockchain) getLastHash() string {
// 	if len(b.blocks) > 0 {
// 		return b.blocks[len(b.blocks)-1].hash
// 	}
// 	return ""
// }

// func (b *blockchain) addBlock(data string) {
// 	newBlock := block{data, "", b.getLastHash()}
// 	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
// 	newBlock.hash = fmt.Sprintf("%x", hash)
// 	b.blocks = append(b.blocks, newBlock)
// }

// func (b *blockchain) listBlocks() {
// 	for _, block := range b.blocks {
// 		fmt.Printf("Data: %s\n", block.data)
// 		fmt.Printf("Hash: %s\n", block.hash)
// 		fmt.Printf("Prev Hash: %s\n", block.prevHash)
// 	}
// }

// B1
// 	b1Hash = (data + "")

// B2
// 	b2Hash = (data + b1Hash)

// B3
// 	b3Hash = (data + b2Hash)
