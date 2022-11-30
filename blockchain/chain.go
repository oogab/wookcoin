package blockchain

import (
	"sync"
)

// 더이상 block slice가 아님, 마지막 블록을 알아야 한다.
// 정확히는 블록체인의 마지막 hash를 알아야 한다.
type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

// 새로운 block을 만들거지만 slice 같은 곳에 추가하지는 않을 거다.
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height
}

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis")
		})
	}
	return b
}
