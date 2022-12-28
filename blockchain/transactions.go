package blockchain

import (
	"time"

	"github.com/oogab/wookcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
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

// 받을 사람 to와 amount를 넘겨준다.
// 보내는 사람은 중요하지 않다. 왜? 보내는 사람은 function에서 받아오는게 아니고 지갑에서 받아올 것이기 때문에
// 만약 어떤 이유로 transaction을 추가할 수 없게되면 error를 발생시킬 수도 있다.
func (m *mempool) AddTx(to string, amount int) error {
}
