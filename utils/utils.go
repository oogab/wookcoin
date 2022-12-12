package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// 이 function에서는 뭐든지 받을 수 있게 하겠다.
// interface는 base type 같아서 뭐든지 interface가 될 수 있다.
func ToBytes(i interface{}) []byte {
	// buffer는 bytes를 넣을 수 있고 read-write 할 수 있다.
	var aBuffer bytes.Buffer
	// block을 Encode한 다음, 그 결과를 aBuffer에 저장
	encoder := gob.NewEncoder(&aBuffer)
	HandleError(encoder.Encode(i))
	return aBuffer.Bytes()
}
