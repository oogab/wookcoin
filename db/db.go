package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/oogab/wookcoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"

	checkpoint = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
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

func SaveBlock(hash string, data []byte) {
	fmt.Printf("Saving block %s\nData: %b\n", hash, data)
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleError(err)
}

// data -> BlockChain NewestHash, Height 이 두개 값만 저장된다.
// 그래서 key 자체는 크게 상관없다. 내가 지어주고 싶은대로 입력하면 됨
func SaveBlockchain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleError(err)
}

// 이 함수는 argument를 받지않고 bytes의 slice를 return
// 이것을 통해서 블록체인에 checkpoint가 있는지 없는지를 알 수 있다.
func Checkpoint() []byte {
	var data []byte
	// DB에 read-only transaction인 View()를 써준다.
	DB().View(func(t *bolt.Tx) error {
		// data bucket을 가져온다.
		bucket := t.Bucket([]byte(dataBucket))
		// bucket에서 checkpoint key로 data가져온다.
		// bucket.Get()은 error를 return해 주지 않는다.
		// byte의 slice만 return 해준다.
		data = bucket.Get([]byte(checkpoint))
		// 그래서 return nil
		return nil
	})
	return data
}
