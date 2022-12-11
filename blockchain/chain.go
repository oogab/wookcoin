package blockchain

import (
	"sync"
)

// 블록체인의 마지막 hash를 알아야 한다.
// 다음 블록의 Height를 알기 위해서, 블록체인이 얼마나 긴지도 알아야 한다.
type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var (
	b    *blockchain
	once sync.Once
)

// 블록을 DB에 저장하겠다.
// 이 함수가 처음으로 호출되면 PrevHash랑 Height가 필요한데 blockchain에 없다.
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = block.Height
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			// 사용자가 처음 오면 빈 내용의 NewestHash랑 0의 Height로 blockchain을 설정한다.
			b = &blockchain{"", 0}
			b.AddBlock("Genesis")
		})
	}
	return b
}
