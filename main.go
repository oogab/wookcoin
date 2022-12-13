package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {
	// 목표는 0이 2개로 시작하는 hash를 찾는 거고 입력값을 바꿀 수는 없다.
	// 입력값을 바꿀 수는 없지만 하나를 더할 수는 있는데 그게 Nonce
	// Ex) hello -> hello1, hello2, hello3, hello4...
	// 보다시피 검증하는 작업은 아주 쉽다.
	// 하지만 컴퓨터가 하는 일은 아주 힘들다.
	// 왜냐하면 hash 함수는 결정론적이고 일방통행이기 때문에,
	// 즉 입력값을 제출할 때 출력값을 전혀 알 수 없다.
	// Ex) hello9를 넣으면 0으로 시작할거야 -> 이렇게 예측할 수가 없다.

	// 이제 개념증명을 만들어 보자.
	// 개념증명은 이 작업증명에서 해쉬 시작에 0이 몇 개 있어야 되는지 확인한다.
	difficulty := 2
	target := strings.Repeat("0", difficulty)
	nonce := 1
	// 입력값 "hello"를 hash해 준다.
	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte("hello"+fmt.Sprint(nonce))))
		fmt.Printf("Hash:%s\nTarget:%s\nNonce:%d\n\n", hash, target, nonce)
		if strings.HasPrefix(hash, target) {
			return
		} else {
			nonce++
		}
	}
}
