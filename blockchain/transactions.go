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
	// 이제 더 이상 Amount는 필요하지 않다.
	// 기존 output에서 받아올 것이기 때문에
	// 기존 Tx의 ID를 TxID에 저장
	// 그리고 그 Tx의 TxOuts들 중 몇 번째 것을 참조할 지 Index 저장
	//* 즉 TxIn은 TxOut을 찾아갈 길과 같은 역할을 한다.
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// 다음으로 진행하기 전 새로운 type struct를 정의한다.
// 우리가 어떤 output이 쓰였는지 안 쓰였는지 확인할 수 있게 도와줄 예정
type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	// 이제 TxIn은 transaction ID, index 그리고 owner가 필요하다.
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
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

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough 돈")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	// transaction output으로부터 받아온 총 잔고가 저장된다.
	total := 0
	uTxOuts := Blockchain().UTxOutsByAddress(from)
	for _, uTxOut := range uTxOuts {
		if total > amount {
			break
		}
		// 이제 amount를 transaction input에 전달해주지 않는다.
		// 하지만 아직 owner를 input에 저장하고 이는 보안상 안전하지 않다.
		// 왜냐하면 사람들이 이 owner를 그냥 변경할 수 있기 때문 -> 이후 지갑 구현에서 보안 적용 예정
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
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

// * 아래 메소드는 '승인할 트랜잭션들 가져오기'같은 느낌의 함수이다.
// 모든 transaction들을 건네주고 mempool 또한 비워줘야 한다.
// 아래 메소드의 로직은 모두 블록을 채굴했을 때에만 실행된다.
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("wook")
	// Mempool에서 모든 tx들을 가져오고 coinbase tx와 함께
	// txs라는 큰 배열에다가 담아준다.
	txs := m.Txs
	txs = append(txs, coinbase)
	// Mempool의 tx들을 비워준다.
	m.Txs = nil
	return txs
}
