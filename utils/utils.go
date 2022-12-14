package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}

func Hash(i interface{}) string {
	// %v는 기본 formatter이다.
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
