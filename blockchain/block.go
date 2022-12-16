package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/oogab/wookcoin/db"
	"github.com/oogab/wookcoin/utils"
)

// 이제 block에서 몇가지 작업을 하자
type Block struct {
	// Data는 단순히 사용자가 입력한 문자열이니 지워주자!
	// Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`
	// 우리의 Block은 이제 transaction을 가지게 된다.
	Transactions []*Tx `json:"transactions"`
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

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)

	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("\n\n\nTarget:%s\nHsah:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

// Data 관련 parameter, field 삭제
func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().difficulty(),
		Nonce:      0,
		// 채굴자들을 위한 거래내역들이 들어갈 예정이다.
		// 이게 코인베이스 거래이다. -> 블록체인에서 생성되는 거래내역
		// 여기에 코인베이스 거래내역을 생성하는 함수가 필요하다.
		Transactions: []*Tx{makeCoinbaseTx("wook")},
	}
	block.mine()
	block.persist()
	return block
}
