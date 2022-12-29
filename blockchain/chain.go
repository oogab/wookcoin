package blockchain

import (
	"sync"

	"github.com/oogab/wookcoin/db"
	"github.com/oogab/wookcoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	// 2분마다 블록 1개 생성
	blockInterval int = 2
	allowedRange  int = 2
)

// 블록체인의 마지막 hash를 알아야 한다.
// 다음 블록의 Height를 알기 위해서, 블록체인이 얼마나 긴지도 알아야 한다.
type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrnetDifficulty int    `json:"currnetDifficulty"`
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

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	// 블록을 추가할 때마다 현재의 difficulty 저장
	b.CurrnetDifficulty = block.Difficulty
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

func (b *blockchain) recalculateDifficulty() int {
	// 5개의 블록을 만드는 데 얼마나 시간이 걸렸는지 파악해야 한다.
	// 예상으로는 블록 1개 생성 시간이 2분이므로 difficultyInterval * blockInterval = 10분
	allBlocks := b.Blocks()
	// 가장 최근 블록 정보를 가져오자.
	// 가장 최근 블록은 slice 맨 앞에 있다. 왜냐하면 가장 최근 hash부터 찾기 때문이다.
	newestBlock := allBlocks[0]
	// 5번째로 최근에 만들어진 블록 정보를 가져오자.
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	// 이제 위에서 구한 두 블록들의 생성 사이에 걸린 시간을 찾아보자.
	// Timestamp는 초 단위이므로 60으로 나눠서 분 단위로 변경한다.
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval

	// 아래 로직은 너무 엄격한 면이 있다
	// 정확히 10분이 아닌 경우 반드시 조건문을 타게 되기 때문에,
	// 어느정도 범위 내에 있다면 조건문을 타지 않게 해준다.
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrnetDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrnetDifficulty - 1
	}
	return b.CurrnetDifficulty
}

func (b *blockchain) difficulty() int {
	// 블록체인의 height가 0이라면, 즉 새 블록체인이라면 2를 return
	// 이 값은 계속 사용되니 defaultDifficulty라는 변수로 지정
	if b.Height == 0 {
		return defaultDifficulty
		// 실제 블록체인은 2016개 마다 difficulty를 조정하지만,
		// 우리는 5개로 한정한다. difficultyInterval = 5
	} else if b.Height%difficultyInterval == 0 {
		// recalculate the difficulty
		return b.recalculateDifficulty()
		// 첫번째 블록이 아니고 Difficulty를 재설정 할 때가 아니라면
	} else {
		return b.CurrnetDifficulty
	}
}

// 여기 이해가 잘 안되는데 직접 block이랑 Ins, Outs를 그리고 추적하면 이해가 빠르다.
func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	// 사용한 Output들 정보를 저장하는 map을 만든다.
	// key -> transaction ID, value -> true, false
	creatorTxs := make(map[string]bool)

	for _, block := range b.Blocks() {
		// 모든 transaction에 대해서 실행
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				// input을 통해서 input을 생성할 때 사용한 output이 포함된 transaction을 찾는다.
				// 그리고 해당 transaction의 output은 이미 사용되었음을 저장
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			// output이 앞서 생성한 creatorIds안에 있는 transaction 내부에 없다는 것을 확인
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					// ok값은 기본적으로 이 key 값이 이 map 안에 있는지 없는지 확인을 해주는 값
					// ok가 true라면 이 output을 생성한 transaction을 creatorTxs에서 찾았다는 것이고,
					// 이 output은 이미 다른 transaction에서 input으로 사용되었다는 증거
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOuts = append(uTxOuts, &UTxOut{tx.Id, index, output.Amount})
					}
				}
			}
		}
	}

	return uTxOuts
}

// 거래 출력값들의 총량을 보여준다.
func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.UTxOutsByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()

			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}
