package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/oogab/wookcoin/db"
	"github.com/oogab/wookcoin/utils"
)

// 0이 2개로 시작하는 hash를 찾는다.
const difficulty int = 2

// hash는 결정론적 함수이다. -> 출력값을 바꾸려면 입력값을 수정해야한다.
// 그런데 블록체인에서는 뭔가를 수정할 수 있는게 거의 없다.
// hash를 수정하면 이 블록을 사용하지 못하게 됨
// Data는 유저가 보내주는 거라 수정하지 못 함
// Height도 마찬가지로 수정하지 못 함
// 하지만 Noncd 값은 블록체인에서 채굴자들이 변경할 수 있는 유일한 값이다.
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
