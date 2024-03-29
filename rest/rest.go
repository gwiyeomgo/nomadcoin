package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gwiyeomgo/nomadcoin/blockchain"
	"github.com/gwiyeomgo/nomadcoin/p2p"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"github.com/gwiyeomgo/nomadcoin/wallet"
	"log"
	"net/http"
)

type url string

//const port string = ":4000"
var port string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
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

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addPeerPayload struct {
	Address, Port string
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	//쿼리 파라미터 받기
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(blockchain.Blockchain(), address)
		json.NewEncoder(rw).Encode(balanceResponse{Address: address, Balance: amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(blockchain.Blockchain(), address)))
	}
}
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Encode 가 Marshal 일을 헤주고
		//결과를 ResponseWriter 에 작성
		//rw.Header().Add("Content-Type", "application/json")
		//json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		//json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		//request 의 body를 받는다.
		//rest client 로 insomnia ,postman 등 사용
		//request body를 struct 로 decode 한다
		//read 할땐 decode
		//& 포인터를 더해주면 addBlockBody 주소를 전달 (복사본 x)
		//*json을 go로 변환시키는 법
		//var addBlockBody AddBlockBody
		//utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		//request body 의 message로 새 블록 추가한다
		//blockchain.Blockchain().AddBlock(addBlockBody.Message)
		newBlock := blockchain.Blockchain().AddBlock()
		//모든 네트워크에 p2p.BroadcastNewBlock() 전달
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated) //201

	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}
func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
	//아래처럼 즉석으로 쓸 수 있음
	//json.NewEncoder(rw).Encode(struct{Address string `json:"address"`}{Address: address})
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	//reqeust body 값을 받고
	json.NewDecoder(r.Body).Decode(&payload)
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	// go 루틴을 사용해서
	// Header 가 broadcast 를 기다리지 않도록 할 수 있다
	go p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//API 에서 온 json 을 go 언어로 variable 바꿀 준비
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
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
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See a block",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an address",
		},
		{
			URL:         url("/mempool;"),
			Method:      "GET",
			Description: "See Mempol",
		},
		{
			URL:         url("/transaction"),
			Method:      "POST",
			Description: "ADD Transaction to Mempool",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to Web Sockets",
		},
		{
			URL:         url("/peers"),
			Method:      "POST",
			Description: "Connect to peer",
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
	//rw.Header().Add("Content-Type", "application/json")
	//fmt.Fprintf(rw,"%s",b)
	//5.더 쉬운 방법 json.NewEncoder()
	// data(struct)을 encode 해서 writer 에 담아 보냄
	json.NewEncoder(rw).Encode(data)

}
func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//height, err := strconv.Atoi(vars["height"])
	//strconv.Atoi string to int
	//utils.HandleErr(err)
	hash := vars["hash"]
	//height는 string 이기때문에 int 로 convert
	//strconv 패키기 이용
	//block, err := blockchain.GetBlockchain().GetBlock(height)
	block, err := blockchain.FindBlock(hash)

	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

//모든 request에 content-type을 설정하는 middlewares 추가하기
//middleware 는 function 인데 먼저 호출되고
//다음 function을 부르고 그럼 거기서 또 다음 function을 부른다

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		//response
		rw.Header().Add("Content-Type", "application/json")
		//(2) 다음 handlerFunc 호출
		next.ServeHTTP(rw, request)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL)

		next.ServeHTTP(rw, request)
	})
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

	//middleware function
	//마지막 목적지 전에 호출되는 녀석
	//(1)호출
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	//(3) HandleFunc 호출
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	//router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	//hexadecimal 을 a-f 와 숫자를 갖는 포맷
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgreade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
