package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/oogab/wookcoin/db"
	"github.com/oogab/wookcoin/utils"
)

// block에 두 가지 요소를 추가한다.
// block 구조를 변경했기 때문에 .db 파일은 삭제한다.
type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
}

// block을 저장하기 위해 만들어 놓은 SaveBlock을 호출한다.
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

// 빈 블록을 만들고 그 블록을 이용하여 복원할 예정
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
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
