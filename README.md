# golang basic

# Variables in Go
```go
var name string = "test"
// type 이 compiler 에게 name 은 항상 string 이라고 말해준다
name1 :="test"
//compiler 가 자동으로 name의 type이 string 이라는 걸 안다.
const name string =""

```

# What is type int in Golang?
int         :signed integer (정수)
,int 뒤의 숫자는 크기를 나타냄
ex> unsigned 8-bit integers (0 to 255)
unit        :unsigned integer(= 음의 정수 아닌 양의 정수만 해당)


[출처](https://go.dev/tour/basics/11)

# Functions
argument
```go
package main
import (
	"fmt"
)
func plus (array ...int)  int{
	// a는 int의 array
	results := 0
	for _,item :=range array {
		results += item
	}
	return results
}
func main()  {
	result := plus(1,1,1,1,1,2)
	fmt.Println(result)
	name :="gwiyeom go~~!!#!#!"
	for _,item :=range name{
		//fmt.Println(item) //byte 타입으로 출력된다.
		fmt.Println(string(item))
	}
}

````

```
multiple return 



# 'text/template' HTML templage 가 구현하는 package
https://pkg.go.dev/text/template

main에서 
pages,partials 폴더에 만든 template을
일일이 load 하지 않고
home 에 대한 요청이 있을 때 마다
home function 이 호출 될 때마다
home template 을 parsing 하고
main 에 load

main 코드에  variable 을 추가했는데
해당 variable은 전체 template을 관리한다

# interface

목표: 특정 string 을 모든 url 앞에 더해준다
https://go.dev/tour/methods/17
interface 는 함수의 청사진과 같다
Stringer 라는 interface 는 String 이라는
하나의 method 만 구현시킨다
대문자로 시작하며 ,매개변수 x, string 을 return

fmt package 로 출력 할때 
어떻게 보여줄지 조절할 수 있다
* Go 에서는 모든 interface 가 내재적으로 구현


* Marshal,Unmarshal 할때, Field 의 결과물을 수정할 수 있는 interface 가 이다
interface 를 json으로 변환할때
중간에서 field 가 json에서
어떻게 보일지 원하는 대로 변경 가능
https://pkg.go.dev/encoding

type TextMarshaler

/*func (u URLDescription) String() string{
	return fmt.Sprintf("http://localhost:4000%s",u.URL)
}*/
//6.Method 대문자는 public
// 소문자로 쓰고 싶을때
// field struct tag 사용
// field struct tag로 json 형태 key 값으로 보내짐
//* omitempty 는 field 가 비어있으면 field 를 숨겨준다.
//* field is ignored by this packages
// `json:"-"` 사용 field 를 무시
//7.MarshalText 는 Field가 json string으로써 어떻게 보일지 결정하는 method


# NewServeMux

```
// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

```
Mux는 multiplexer 
url로 request 를 다루는 것
url을 지켜보고 내가 원하는 함수를 실행
client 가 request 를
보내면 multiplex가
그 request 를 확인

Gorilla mux
go get -u github.com/gorilla/mux

# Atoi

# adapter 패턴

http.HandlerFunc
사실 이건 function이 아니라 type이다
http.HandlerFunc 이 return 하는 http.Handler
는 interface
이 interface는 ServeHTTP 라는 method 를 구현

HandlerFunc 라는type은 바로 adapter 이다
HandlerFunc 은 매개변수로
func(rw http.ResponseWriter,r *http.Request){
}) 
를 쓴다면
함수를 호출하는 것이 아닌 type을 만드는 것이다

HandlerFunc type 이 어떻게
Handler interface 로 인식되는가?
adapter 의 힘
adapter 에게 적절한 argument 를 보내주면
그 다음, adapter 는 네가 필요한 것들을 구현해준다.

type? 
우리가 보낸 function이 이 조건에 부합하는지를 보고
receiver function의 훌륭함 덕택에
Handler 가 가져야 하는 ServerHTTP method 를 
구현할 거다
즉, 우리가 직접 struct 혹은 type을 만들어서ServeHTTP 를 구현하는 대신
adapter가 하는 것은
여기 명시된 형태(type 의 유형)에 맞는 겋 보내라고 알려준다

# CLI
Command Line Interface

standard library 에 있는 flag 사용
CLI를 구축하는 것을 도와준다

참고로
go 로 CLI를 만드는 framework 가 있는데
cobra 이다
cobra 는 CLI를 만들 때 필요한 많은 도우미 function 을 제공해준다.


1. 일단 console과 상호작용하면서
유저가 보낸 command 는 어떻게 얻을 수 있는지
os.Args
```
C:\github\golang\nomadcoin>go run main.go rest
[C:\Users\me\AppData\Local\Temp\go-build3973503983\b001\exe\main.exe rest]
```
#8 Bolt
persistence 는 기본적으로 DB 백엔드를 블록체인에 더한다는 의미
서버를 재시작하면 저장해던 내용이 사라진다
기존 코드에서는 slice block 을 담고
메모리에 저장했었다.
block slice 는 memory 에 저장할 필요가 없다
db에 검색하는 기능

go의 key/value db 인 bolt
bolt 를 선택한 이유?
bolt 프로그램은 안정적이라 더 바뀔 일 없기 때문에
완전히 완성되었다 봅니다.
더 이상 변화 없음

변화무쌍한 자바스크립트 생태계와 다르게..

* DB의 경우
1. DB파일이 존재하지 않으면 파일을 만들어서 initialize 한다
`go get github.com/boltdb/bolt/...`


# boltbrowser
comand line interface로
bolt 파일을 열어서 확인해 볼 수 있다.
https://github.com/br0xen/boltbrowser
`go get github.com/br0xen/boltbrowser`

`boltbrowser <filename>` => bucket 생성

creae root bucket 에 blocks 를 써주고

go build main.go 실행
`boltbrowser <filename>` 실행해더니
+data가 생겼음

-> 확인 못합
https://github.com/evnix/boltdbweb
go get github.com/evnix/boltdbweb


https://stackoverflow.com/questions/41836209/only-one-usage-of-each-socket-address-protocol-network-address-port-is-normall

Troubleshooting Port Conflicts (Only one usage of each socket address is normally permitted)
=> port conflict
https://kb.fastvue.co/fortigate/s/article/troubleshooting-port-conflicts-only-one-usage-of-each-socket-address-is-normally-permitted
https://m.blog.naver.com/PostView.naver?isHttpsRedirect=true&blogId=ysw1130&logNo=220159168596
port 찾기
netstat -ano | find "8080"

This will list all the processes on the machine using port 8080 (it may also include other processes that have a substring of 8080).

1.어플리케이션을 처음실행
메모리에 저장하는 것이 아니라
필요하면 DB를 찾아서 화면에 보여지도록 작업
Blockchain() 실행시 
블록체인이 initialize 되고 (singleton)

빈 블록체인 만들고
DB에서 블록체인에 checkpoint 가 저장되어 있는지 확인한다.

JSON을 encoding/decoding 작업
gob package를 이용
//ToBytes
//FromBytes

chckpoint 를 추가해서 block을 저장

2.db에 블록이 저장되어 있을때
blockchain.db에서 checkpoint를 불러와서 byte를 찾음
블록체인 복원
->FindBlock 은 hash를 받고 Block 포인터 반환

Get All Block
Blocks 함수




