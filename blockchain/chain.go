package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
	decoder := gob.NewDecoder(bytes.NewReader(data))
	// Decode는 pointer로 decode해야 한다고 되어 있다.
	// pointer가 아닌 것을 보내면 decode할 수 없다.
	// decode는 memory address의 value를 변경하기 때문에.
	// blockchain 전체를 이 decode된 값으로 바꾼다.
	utils.HandleError(decoder.Decode(b))
	// DB에서 찾은 byte를 텅 빈 블록체인의 memory address에 decode한다.
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

// 블록을 DB에 저장하겠다.
// 이 함수가 처음으로 호출되면 PrevHash랑 Height가 필요한데 blockchain에 없다.
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
			checkpoint := db.Checkpoint()
			// checkpoint가 없다면 연결된 block이 하나도 없다는 것
			// Genesis block을 만든다.
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				// restore blockchain from bytes
				// bytes를 decode한다.
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
	return b
}
