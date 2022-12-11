package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/oogab/wookcoin/db"
	"github.com/oogab/wookcoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

// block을 byte slice로 바꿔주는 함수
// gob package를 사용, bytes를 encode/decode 하게 해준다.
func (b *Block) toBytes() []byte {
	// buffer는 bytes를 넣을 수 있고 read-write 할 수 있다.
	var blockBuffer bytes.Buffer
	// block을 Encode한 다음, 그 결과를 blockBuffer에 저장
	encoder := gob.NewEncoder(&blockBuffer)
	utils.HandleError(encoder.Encode(b))
	return blockBuffer.Bytes()
}

// block을 저장하기 위해 만들어 놓은 SaveBlock을 호출한다.
func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}
