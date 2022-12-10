package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oogab/wookcoin/blockchain"
	"github.com/oogab/wookcoin/utils"
)

// const port string = ":4000"
var port string

// URL -> url
type url string

// URL이 어떻게 json으로 Marshal 될지를 정할 수 있다.
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// URL이 대문자일 필요는 없음
// export할 건 Start() 하나 뿐임
// URLDescription -> urlDescription
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

// AddBlockBody -> addBlockBody
type addBlockBody struct {
	Message string
}

func (u urlDescription) String() string {
	return "Hello I'm the URL Description"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{id}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		// Encode가 Marshal의 일을 해주고, 결과를 ResponseWrite에 작성.
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		// user가 {"data": "my block data"}의 형태로 데이터를 보내줌.
		var addBlockBody addBlockBody
		// Decode에 pointer를 보내야 한다.
		// r.Body를 Decode한 후 결과를 addBlockBody에 저장
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	// 새로운 ServeMux를 생성
	// ServeMux는 url(/blocks)과 url 함수(blocks)를 연결해주는 역할을 한다.
	handler := http.NewServeMux()
	port = fmt.Sprintf(":%d", aPort)
	handler.HandleFunc("/", documentation)
	handler.HandleFunc("/blocks", blocks)
	fmt.Printf("Listening on http://localhost%s\n", port)
	// handler에 nil이 아니라 우리가 만든 handler(ServeMux)를 넣어준다.
	// DefaultServeMux를 사용하지 않게 된다.
	log.Fatal(http.ListenAndServe(port, handler))
}
