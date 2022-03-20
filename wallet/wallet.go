package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/utils"
)

func Start() {
	//비공개키 생성
	//go의 표준 라이브러리에 포함된 패키지 사용
	//Elliptic Curve Digital Signature Algorithem
	//rand?
	//이것은 암호화적으로 보안된 난수ㅐㅇ성기의 전역 공유 인스턴스
	//키생성에 난수(randomness) 가 필요하기 때문
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	//매우 큰 숫자인 우리의 피공개키 출력
	//fmt.Println("Private Key",privateKey.D)
	//privatekey 는 publickey 필트를 갖고 있고 publickey를 바로 참조할 수 있다
	//fmt.Println("Public Key,x,y",privateKey.X,privateKey.Y)
	message := "I love you"
	hashMessage := utils.Hash(message)
	// 이 hash 에 서명 (sign) 한다
	//hashMessage 는 string 이니까 []byte 로 바꿈
	//1. []byte(hashMessage)
	//2 .
	hashAsByte, err := hex.DecodeString(hashMessage)
	//hash가 올바른 16 진수 형식을 가지고 있지 않는다면 err 발생
	utils.HandleErr(err)
	//Sign 함수를 써서 난수생성기,비공개키,hash된 byte 넘김
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsByte)
	utils.HandleErr(err)
	fmt.Printf("R:%d\nS:%d", r, s)
	//서명은 R 과 S 값으로 나눠져 있다.
	//r 과 s 통합 방법?
}
