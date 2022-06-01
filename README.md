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

#9 블록체인에 채굴에 대한 작업증명 pow
https://www.youtube.com/watch?v=ElGBP90XZWE

블록체인은 데이터베이스
데이터가 블록에 있고
블록들이 연결되어 있다.

누가 블록을 추가할 수 있지?

작업증명 => 블록체인을 보호함

마이너? 채굴자?
채굴자는 블록체인에 들어오는 데이터를 확인한다
채굴자가 데이터를 블록 안에 넣어서
블록체인에 보내는 역할을 한다.

ex)
채굴자가 해당 거래내역을 확인하고
비트코인 블록체인에서 빗코를 친구에게 보냄
체크가 끝나면 데이터를 블록안에 넣는다.
체크를 여러개 해서 블록이 꽉 차면
블록을 닫고 블록체인에 올림
=>탈중화방법
누구나 원한다면 채굴자가 될 수 있다.

채굴자들이 트랜잭션 컴펌하면 수수료로 돈 받음
검증 작업이 활발하게 이루어짐 but 쉽지x

작업증명이 채굴자에 질문을 한다
채굴자는 질문에 답해야 한다
딥을 찾아야 블록을 올릴 수 있다.

비트코인은 사람들이 블록을 블록체인에 올리는 순간 생성
coinbase transaction => 비트코인 생성 순간

어떻게 컴퓨터가 블록을 더하기를 어렵게 할까?
컴퓨터가 풀기는 어렵지만 검증하기 쉬운 방법 필요
강의에서 방법
작업증명-> n개의 0으로 시작하는 hash 찾기

비트코인에서도 네트워크가 보유한 파워에 따라
difficulty가 변경된다.
nonce 값은 채굴자들이 변경할 수 있는 유일한 값


# transaction

거래?
TxIn[]
TxOut[]
입력값은 너의 주머니의 돈
거래 출력값은 그 거래가 끝났을 때 각각의 사람들이 갖고있는 액수

eX)
내가 너에게 5딸라를 보내고 싶다
$5인지 확인

Tx[$5(blockchain)]
TxOut[$5(miner)]

코인베이스 거래?
블록체인에서 생성되는 거래내역

메모리풀?
Mempoll 은 아직 확정되지 않은
거래 내역들을 두는 곳
미확인 거래내역들로 이루어진 array or slice

채굴자들이 블록을 채굴한 다음에
엄청 비싼 전력량과 컴퓨팅 파워랑
이것저것을 다 준 다음에
mempool에 와서 이 거래내역들을 블록에 추가
mempool 은 아직 확인되지 않은 거래내역들을 보관하는 곳

# uTxOuts

Tx1 
	TxIn[COINBASE]
	TxOut[$5(you)] (1) Spent TxOut
Tx2 
	TxIn[Tx1.TxOuts[0]] (2) 은 (1)이전 의 트랜잭션 output 과 연결되어야 한다
=> 이방법은 trnsacion의 Id를 검색
	TxOuts[$5(me)] 
Tx3
 	TxIns[Tx2.TxOuts[0]]
	TxOuts[$3(you), $2(me)] => unspend transacion output * 2

# transacion을 mempool 에 올리는 것을 제한

함수가 receiver function 혹은 method 여야 하는지 아닌지 알려주는 규치?
object -oriented programming 에서
method란 클래스 내부에 존재하는 함수
go 에서는 클래스는 없고
구조체 struct 만 존재한다

함수(function)이 구조체(struct) 를 
변화시킨다면 그 함수는 메서드 여야 한다
하지만 sruct 가 변화하지 않는다면
그건 메서드가 아니다

#wallet
1.gwiyeom 이 unspent transaction output 을 소유하고 있는지 확인
2.gwiyeom 이 트랜잭션을 승인했는지 검증
=>서명 signiture
트랜잭션이 gwiyeom 에 의해서
바랭하고 승인됐다는 것을 확인가능

1, 서명 dignature ,검증 verification등
어떻게 동작하는지 확인
공개키 public key vs 비공개 private key 암호화
2. 지감 유지 persistence 영속성
지갑 파일로 저장,복구 방법
3. 서명,증명을 구현


# 서명된 메세지를 보내는 방법
1. we have the message
"message" -> hash (x) => "hashed_message"

2.generate key pair
KeyPair (privateK, publick) (save priv to a  file)
// 비공개키를 파일로 저장
//비공개키가 남아있지 않다면 (잃어버리면) 서명할 수 없다

3. sign the hash
("hashed_message"+privateK) -> "signature"
비공개키는 노출되면 안된다
노출된다면 누군가 나인척하고 서명할 수 있다
T
비공개키로 서명하고
공개키로 검증을 한다.
4. verify
("hashed_message"+"signiture"+publick) -> true/false
세 가지의 조합으로
해당 비공개키로 이 메세지가 서명되었는지 검증한다


# wallet 은 privare Key 와 address 를 가지고 있다
public key 는 공유되고
private key 는 바깥 세상의 그 누구와도 공유되지 않는다

싱글톤 패턴으로 wallet 파일이 이미 존재하는지 확인하는데
restoreKey 함수에서는
파일을 읽고
x.509 패키지를 사용해서 private key 를 복구
wallet 없는경우는 createPrivKey 를 해서 privateKet를 생성한다

privateKey 로 address 를 얻음
16진수 문자열로 바꿔서 return 한것이 address

//11.11 Transaction Signing

#11.14
2개의 transaction output 을 갖고 있음
(TxOut1,TxOut2)

Tx
	TxIn[
		(TxOut1)
		(TxOUt2)
	]
	Sign : X
	
	transaction 을 만든사람이output 을 갖고있음
	
//private key 로 서명하고
//public key 로 검증할 수 있다
TxIn.Sign + TxOut1.Address = true /false



