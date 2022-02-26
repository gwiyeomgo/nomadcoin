package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//6.Method 대문자는 public
// 소문자로 쓰고 싶을때
// field struct tag 사용
// field struct tag로 json 형태 key 값으로 보내짐
//* omitempty 는 field 가 비어있으면 field 를 숨겨준다.
//* field is ignored by this packages
// `json:"-"` 사용 field 를 무시
type URLDescription struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

const port string = ":4000"

func documentation(rw http.ResponseWriter, r *http.Request) {
	// data 는 Go의 세계에 있는 slice
	//struct 의 slice
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "test",
		},
		{
			URL:         "/blocks",
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
func main() {

	http.HandleFunc("/", documentation)
	//REST API를 통해서 transaction 만들기
	//1.서버를 시작'
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
	//explorer.Start()
}
