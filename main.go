package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data    string
	hash    string
	preHash string
}
type blockchain struct {
	blocks []block
}

func main() {
	chain := blockchain{}
	chain.addBlock("B block")
	chain.addBlock("S block")
	chain.addBlock("T block")
	chain.listBlocks()
}

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
}

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
