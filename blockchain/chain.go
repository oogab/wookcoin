package blockchain

import (
	"sync"

	"github.com/oogab/wookcoin/db"
	"github.com/oogab/wookcoin/utils"
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

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

// 블록을 DB에 저장하겠다.
// 이 함수가 처음으로 호출되면 PrevHash랑 Height가 필요한데 blockchain에 없다.
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

// NewestHash를 가지고 해당 블록을 찾는다.
// 그리고 가장 최근 블록의 prevHash로 이 과정을 반복한다.
// prevHash가 없는 블록을 찾을 때 까지 계속
func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			// prevHash가 없는 Genesis 블록에 도달했다면! for문 탈출
			break
		}
	}
	return blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			// checkpoint가 없다면 연결된 block이 하나도 없다는 것
			// Genesis block을 만든다.
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				// restore blockchain from bytes
				// bytes를 decode한다.
				b.restore(checkpoint)
			}
		})
	}
	return b
}
