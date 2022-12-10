package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func block(rw http.ResponseWriter, r *http.Request) {
	// 이제 함수에 필요한 이 id를 어떻게 받아오는지 알아보자.
	vars := mux.Vars(r)
	id := vars["id"]
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	// 이 router의 특징은 url이 어떤 method를 처리할지 특정할 수 있다는 것
	// 이게 유용한 이유는 다른 method로부터 보호해 주기 때문이다.
	// 만약 이게 가능하지 않았다면 우리는 router function 내부의 switch case에서 default 처리가 필요하다.
	// default:
	//   rw.WriteHeader(http.StatusMethodNotAllowed)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{id:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
