package main

import (
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/blockchain"
	"html/template"
	"log"
	"net/http"
)

const port string = ":4000"

//struct 내부값을 공유하기 위해 대문자로 쓴다
type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
	//blockchain package 의 pointer slice
}

//Go 로 Server-side 렌더링 웹사이트 만들기
//하나의 library 만 있으면 돼
func home(writer http.ResponseWriter, request *http.Request) {
	// http.ResponseWriter : 유저에게 보내고 싶은 데이터
	// http.Request : 유저가 보낸 데이터
	//console 에 보여주지 않고 writer 에 출력
	//즉 data 를 format 해서 writer 에 보낸다.
	//fmt.Fprint(writer,"응답값")
	//HTML template 을 렌터링
	//template utility
	//ParseFiles  error 발생 시 Must 가 자동으로 error  출력
	//GoHTML is an HTML formatter for Go. You can format HTML source codes by using this package.
	//https://andybrewer.github.io/mvp/ 로 html 파일 꾸미기
	//tmpl := template.Must(template.ParseFiles("templates/home.gohtml"))
	/*tmpl, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		log.Fatal(err)
	}*/
	//go package 중에서 template 들을 load 해오자
	tmpl := template.Must(template.ParseFiles("templates/home.gohtml"))

	data := homeData{
		PageTitle: "HOME",
		Blocks:    blockchain.GetBlockchain().AllBlocks(),
		//blockchain 에 있는 모든 block 을 갖다주는 function
	}
	//template으로 data 를 보냄
	//template 내부에서는 pageTitle 이라고 하는 field 를 기다리고 있음
	tmpl.Execute(writer, data)
}
func main() {
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", port)
	//서버 만들기 http://golang.site/go/article/111-%EA%B0%84%EB%8B%A8%ED%95%9C-%EC%9B%B9-%EC%84%9C%EB%B2%84-HTTP-%EC%84%9C%EB%B2%84
	//서버 생성
	//http.ListenAndServe(port,nil)
	//Fatal 은 os.Exit(1)다음에 따라나오는
	//os.Exit(1) 은 프로그램이 error code 1 으로 종료되는 것
	//error 를 print() 하는 것과 동일
	//ListenAndServe function은 error 를 반환한다
	//log.Fatal 은 error 가 있다면 출력
	//없다면 ListenAndServe 은 절대 끝나지 않고,
	//Fatal 은 절대 실행되지 않음
	log.Fatal(http.ListenAndServe(port, nil))
}
