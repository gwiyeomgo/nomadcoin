package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data    string
	Hash    string
	PreHash string
}
type blockchain struct {
	blocks []*Block // pointer 들의 slice
}

//blockchain 을 공유하고 초기화하는 부분을 구현
//singleton 패턴
//우리의 application 내에서 언제든지
//blockchain 의 단 하나의 instance 만을 공유하는 방법

//변수 b 선언, 타입은 blockchain 으로 pointer type
//변수를 소문자로 써서 private (blockchain 패키지 내에서만 사용)
//singleton 의 의미는,이 변수의 instance 를 직접 공유하지 않고
//이 변수의 instance 를 우릴 대신해서 드러내주는 function 을 생성하는 것
var b *blockchain

//다른 패키지에서 blockchain이 어떻게 드러날 지를 제어할 수 있다는 의미
//다른 곳에서 blockchain을 요청하면 초기화 후 반환

//sync 패키지
//우리가 동기적으로 처리해야 하는 부분을 제대로 처리하게 도와준다
//sync 의 Once 는 Do 라는 function 갖고 있다.
//Once 는 단 한번만 호출되도록 해주는 함수
var once sync.Once

func GetBlockchain() *blockchain {
	if b == nil {
		//b를 초기화 (처음이자 마지막)
		once.Do(func() {
			//우리가 몇 천 개의 goroutine을 실행해도 오직 한번만 호출
			b = &blockchain{}
			b.AddBlock("Abc")
		})
	}
	return b
}
func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func createBlock(data string) *Block {
	newBlock := Block{
		Data:    data,
		Hash:    "",
		PreHash: getLastHash(),
	}
	newBlock.calculateHash()
	return &newBlock
}
func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}
func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PreHash))
	//hash 를 "%x"를 통해 16진수로 변경
	b.Hash = fmt.Sprintf("%x", hash)
}
func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}

/*//해당 function 이 하나의 일만 하도록 refactoring
//복사본이 아닌 원본에 영향이 가길 원하기 때문에 * 포인터를 리시버에쓴다
func (b *blockchain) addBlock(data string) {
	//첫블록을 제외하고 모든 block 은 이전 hash 가 존재
	newBlock := block{
		data:    data,
		hash:    "",
		preHash: b.getLashHash(),
	}
	//hash 생성
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.preHash))
	//hash 를 "%x"를 통해 16진수로 변경
	newBlock.hash = fmt.Sprintf("%x", hash)
	//append 이용 b.blocks 에 newBlock 추가
	//append 는 추가하고픈 element 와 함께 새로운 slice를 반환
	b.blocks = append(b.blocks, newBlock)

}

func (b *blockchain) getLashHash() string {
	length := len(b.blocks)
	if length > 0 {
		// 배열의 마지막 값 = 배얄[배열길이의 -1]
		return b.blocks[length-1].hash
	}
	return ""
}
func (b *blockchain) listBlocks() {
	for _, block := range b.blocks {
		fmt.Printf("Data:%s\n", block.data)
		fmt.Printf("Hash:%s\n", block.hash)
		fmt.Printf("preHash:%s\n\n", block.preHash)
	}
}*/

// one -way function : 단방향으로만 실행
//hash function `h_fn`에 x 매개변수로 값을 넘겨주면 => 오른쪽의 값 출력
//* 같은 입력값은 항상 같은 출력값을 얻는다
//"test" 는 입력값
//"100" 는 출력값
//"test" = h_fn(x) => "100"
//func main() {
//1.block 을 생성하고 hash 하는 방법
// genesisBlock 에 block 값을 초기화 시켜준다
//	genesisBlock := block{data: "genesis Block"}
//genesisBlock.hash = fn(genesisBlock.data +genesisBlock.preHash)
//암호화폐들이 이용하는 hash 방법은 SHA256 알고리즘을 사용
//golang 에는 sha256.Sum256 있음
// Sum256 매개변수는 [32]byte
//* byte의 slice 와 string 차이점?
//string 도 byte 의 slice 일종이기는 함
// 그개서 string을 for loop 로 반복하면 분리된 byte 를 볼 수 있다.

//for  _ , aByte := range "genesis Block" {
//fmt.Printf(aByte) 오류 : fmt.Printf(value of type string)
//golang 에서 string 의 값은 불변이다.
//이 의미는 string(ex> "test") 값을 우리가 변경 할 수 없다는 뜻
//더크게 더작게 만들 수 없다
//또한 byte 의 배열(고정된 길이를 갖은)과 비슷
//fmt.Printf("%b\n",aByte)
//2진법으로 보이는 string 이 분리된 byte
//}
//배열은 고정된 값을 갖기 때문에 Sum256 함수 내에서
//넘어온 byte 를 변경한다는 의미
//byte slice array 를 결과값으로 받음
//	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.preHash))
//fmt.Sprintf("%x",hash) 의 결과값은 string
// hash 를 "%x"를 통해 16진수로 변경
//	genesisBlock.hash = fmt.Sprintf("%x", hash)
//block 들이 이전 block들의 hash를 가리키고 있다

//}
