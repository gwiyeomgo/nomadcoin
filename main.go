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

// one -way function : 단방향으로만 실행
//hash function `h_fn`에 x 매개변수로 값을 넘겨주면 => 오른쪽의 값 출력
//* 같은 입력값은 항상 같은 출력값을 얻는다
//"test" 는 입력값
//"100" 는 출력값
//"test" = h_fn(x) => "100"
func main() {
	//1.block 을 생성하고 hash 하는 방법
	// genesisBlock 에 block 값을 초기화 시켜준다
	genesisBlock := block{data: "genesis Block"}
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
	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.preHash))
	//fmt.Sprintf("%x",hash) 의 결과값은 string
	// hash 를 "%x"를 통해 16진수로 변경
	genesisBlock.data = fmt.Sprintf("%x", hash)
}
