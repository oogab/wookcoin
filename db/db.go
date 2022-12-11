package db

import (
	"github.com/boltdb/bolt"
	"github.com/oogab/wookcoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

// blockchain이 호출할 많은 함수들을 모아둔 파일이 될 예정
// db.go와 상호작용할 파일은 blockchain.go 밖에 없다.

// 먼저 DB파일이 존재하지 않으면 파일을 만들어서 initialize한다.
var db *bolt.DB

// GetBlockChain()에서 한 singleton 패턴과 매우 유사하다.
func DB() *bolt.DB {
	// bolt에는 bucket이 있다. 다른 RDB의 table과 같다.
	// 두 개의 bucket이 필요
	// 첫 번째 bucket에는 모든 블록을 저장한다.
	// 두 번째 bucket에는 블록체인 자체에 대한 정보를 저장한다.
	// 예를들어 lastHash가 뭐였는지를 알려주는 bucket, 왜?

	// bolt는 key/value로 데이터를 저장하기 때문에 데이터 저장 순서를 보장하지 않는다.
	// 그래서 이를 정렬할 수 있는 방법이 필요하다.
	// "dfasdfasdfa":453423452345
	// "dfjiurowikldf":89187938179283
	// "lskdjlfaisdjkf":98172938719280
	// 마지막이 뭐였는지 알면 다음게 뭔지도 알 수 있으니까!
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0o600, nil)
		db = dbPointer
		utils.HandleError(err)
		// bucket이 있는지 확인하고 없으면 만들도록 read-write transaction이 필요하다.
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleError(err)
	}
	return db
}
