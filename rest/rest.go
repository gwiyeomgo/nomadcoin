package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gwiyeomgo/nomadcoin/blockchain"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"log"
	"net/http"
)

type url string

//const port string = ":4000"
var port string

func (u url) MarshalText() (text []byte, err error) {
	url := fmt.Sprintf("http://localhost%s", port, u)
	return []byte(url), err
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type AddBlockBody struct {
	Message string
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Encode 가 Marshal 일을 헤주고
		//결과를 ResponseWriter 에 작성
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		//request 의 body를 받는다.
		//rest client 로 insomnia ,postman 등 사용
		//request body를 struct 로 decode 한다
		var addBlockBody AddBlockBody
		//read 할땐 decode
		//& 포인터를 더해주면 addBlockBody 주소를 전달 (복사본 x)
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		//request body 의 message로 새 블록 추가한다
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated) //201

	}
}
func documentation(rw http.ResponseWriter, r *http.Request) {
	// data 는 Go의 세계에 있는 slice
	//struct 의 slice
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "test",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "data:string",
		},
	}
	//data 를 json 으로 변경
	//Marshal 은 interface => json으로 변환
	//Unmarshal 은 json => object 변환
	//Marshal 은 메모리형식으로 저장된 객체를
	//저장/송신 할 수 있도록 변환해 준다
	//b ,err:= json.Marshal(data)
	//utils.HandleErr(err)
	//3. byte => string
	//(1)string := string(b)
	//(2)string1 := fmt.Sprint(b)
	//(3) fmt.Printf("%s",b)
	//4.json 을 user 에 return

	// go가 struct 를 json 으로 바꾸는 방식
	//fmt.Fprintf => console 아닌 writer 에 작성하고 싶을때
	//string 형태로 writer 에 담아서 보냄
	//4.이때 content-type =json으로 보내기 위해서
	rw.Header().Add("Content-Type", "application/json")
	//fmt.Fprintf(rw,"%s",b)
	//5.더 쉬운 방법 json.NewEncoder()
	// data(struct)을 encode 해서 writer 에 담아 보냄
	json.NewEncoder(rw).Encode(data)

}
func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

}

func Start(aPort int) {
	//go get -u github.com/gorilla/mux 을 통해 gorilla mux 사용
	router := mux.NewRouter()
	//ServeMux 는 url 과 url 함수를 연결해주는 역할
	//handler := http.NewServeMux()
	port = fmt.Sprintf(":%d", aPort)
	//하나의 multiplexer 가 두곳에서 사용되서 에러남
	//port 가 달라도 상관 없음
	//http.HandleFunc("/", documentation)
	//이렇게 바꿈으로
	//ListenAndServe 함수에 기본 multiplexer 가 기본이 아닌 handler 사용하도록
	//.Methods("GET") 을 쓰면 다른 method로 부터 보호해 준다.
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "PUT")
	router.HandleFunc("/blocks/{id:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}