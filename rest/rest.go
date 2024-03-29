package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oogab/wookcoin/blockchain"
	"github.com/oogab/wookcoin/utils"
	"github.com/oogab/wookcoin/wallet"
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

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
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

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
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
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	// 이제 함수에 필요한 이 id를 어떻게 받아오는지 알아보자.
	vars := mux.Vars(r)
	hash := vars["hash"]
	// hash를 받는 함수가 DB에서 블록을 찾고 복원한 뒤 전달
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		// 여기서는 error 메세지를 encode
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

// 모든 request에 Content-Type을 설정하는 middleware를 만들자!
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	// 중요한 점은 아래 http.HandleFunc는 handler가 아니라 type이다.
	// type HandlerFunc func(ResponseWriter, *Request)
	// HandlerFunc라는 type은 바로 adapter이다.
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

// 이제 transaction을 생성하면 여기서 확인할 수 있다.
func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	// request body의 데이터를 decode해서 json -> struct로 변환
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{"not enough funds"})
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	wallet.Start()
	log.Fatal(http.ListenAndServe(port, router))
}
