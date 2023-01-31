package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/oogab/wookcoin/utils"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)

	message := "i love you"
	hashedMessage := utils.Hash(message)
	fmt.Println(hashedMessage)
	hashAsBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	utils.HandleError(err)

	// Sign 함수는 비공개키가 필요하고 Verify 함수는 공개키가 필요하다.
	ok := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)
	fmt.Println(ok)
}
