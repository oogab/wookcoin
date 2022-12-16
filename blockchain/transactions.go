package blockchain

import (
	"time"

	"github.com/oogab/wookcoin/utils"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

func makeCoinbaseTx(address string) *Tx {
	// 거래 입력값은 소유주가 있다. -> COINBASE
	// 거래 입력값에는 총량도 있다. -> 채굴자에게 지급할 액수의 총량
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	// 거래 출력값의 소유주는 채굴자의 주소
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		// tx를 hash해서 Id로 지정해줘야 한다.
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
