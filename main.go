package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oogab/wookcoin/utils"
)

const port string = ":4000"

/**
/
GET
See Documentation

/blocks
POST
Create a new block
*/

type URLDescription struct {
	URL         string
	Method      string
	Description string
}

// 먼저 user에게 JSON을 보내는 것 부터 시작!
// GO에서 뭔가를 받아서 유효한 JSON으로 변환한다는 것
func documentation(rw http.ResponseWriter, r *http.Request) {
	// URLDescription slice 생성
	// 이 data는 Go의 세계에 있는 slice, struct의 slice -> JSON으로 바꿔야 함.
	// Marshal을 사용해야 한다!
	// Marshal은 JSON으로 encoding한 interface(v)를 return 한다.
	// Marshal은 메모리형식으로 저장된 객체를, 저장/송신할 수 있도록 변환해 주는 것
	// Marshal은 Go에서 interface를 받아서 JSON으로 바꿔주는 것
	// Unmarshal은 반대로 JSON을 받아서 Go의 물건으로 바꿔주는 것
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
	}
	// []byte, error return
	// byte slice가 있음, 이 byte는 JSON으로 바뀐 우리의 data
	b, err := json.Marshal(data)
	utils.HandleError(err)
	fmt.Printf("%s", b)
	fmt.Println(r.Header)
}

func main() {
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
