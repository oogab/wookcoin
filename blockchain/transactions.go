package blockchain

import (
	"errors"
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

// transaction 생성
// transaction은 들어오는 input과 같은 금액의 output이 필요하다.
// wook -> input (100)
// output => wook (30), user (70) 이런식으로..
func makeTx(from, to string, amount int) (*Tx, error) {
	// wook 계좌에 충분한 잔금이 남아있는지 확인
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	// transaction Inputs를 생성 = 내가 보내고자 하는 금액과 똑같은 액수
	// 하지만 tx Inputs를 생성하려면 이전 outputs로 부터 생성해야 한다.
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		// 충분한 tx Inputs를 가지고 있을 때 반복문 탈출
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}
	// 잔돈 tx Output 생성
	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	// user에게 건넬 amount만큼의 tx Output 생성
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

// 받을 사람 to와 amount를 넘겨준다.
// 보내는 사람은 중요하지 않다. 왜? 보내는 사람은 function에서 받아오는게 아니고 지갑에서 받아올 것이기 때문에
// 만약 어떤 이유로 transaction을 추가할 수 없게되면 error를 발생시킬 수도 있다.
// mempool에 transaction을 추가만 할 뿐 생성하지는 않는다.
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("wook", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}
