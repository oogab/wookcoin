package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oogab/wookcoin/blockchain"
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

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
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
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
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

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	// node에서 middleware를 사용하는 방법과 매우 유사하다.
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
