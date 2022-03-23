package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"math/big"
)

const (
	signature   string = "c48ce92a8d8ecc189cbf9c98d254cba1b3bab8f63885f7c2884923d672227d6110cefb3932053da20c490a8eb65ef1d8091016b256515f5429ea4d356a8927fa"
	privateKey  string = "3077020101042062146a70b15e2477c5fbbaf22752806797169649d294b72a80dc18d7620f03f0a00a06082a8648ce3d030107a1440342000435fe9e9f64bcff2a3a0297f1b173143ce103e0649c854f2b43115e1ea849ff52a8bf333a454558824e3815e17bfabc5e7a28ae205a1b3b6fd9ef074c80221d66"
	hashMessage string = "c33084feaa65adbbbebd0c9bf292a26ffc6dea97b170d88e501ab4865591aafd"
)

//* const 로 지정된 문자열들로 부터 복구 작업
func Start() {
	//signature string => 2개 32 bytes slice -> slice 반으로 쪼개 big.Int 갑으로 변환
	// 시그니처의 인코딩 박식이 16진수 형식인지 아닌지 확인
	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]
	//fmt.Printf("%d\n\n%d\n\n%d\n\n",sigBytes,rBytes,sBytes)
	//시그니처 복구
	// var bigR, bigS *big.Int 이경우 선언만 한것 (big.Int 의 포인터)
	//초기화 해주기
	var bigR, bigS = big.Int{}, big.Int{}
	//big.Int struct 의 SetBytes 라는 메서드 가 big.Int 에 byte 값 넘겨준다(전달).
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	//privateKey 는 복구 과정 필요
	//비공개키를 받아서 bytes를 반환해주는 MarshalECPrivateKey
	//x509.MarshalECPrivateKey()
	//bytes 를 받아서 비공개키 반환
	//x509.ParseECPrivateKey([]byte(privateKey)) 이렇게 쓰면 privateKey가 실제 16진수 문자열인지 확인도 안하니까 허술
	//먼저 privateKey 문자열의 인코딩 체크
	//비공개키의 인코딩 방식이 16진수 형식인지 아닌지 확인
	privBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)

	private, err := x509.ParseECPrivateKey(privBytes)
	utils.HandleErr(err)
	fmt.Println(private)
	//hashMessage 는 hash => bytes
	hashBytes, err := hex.DecodeString(hashMessage)

	//ecdsa 패키지의 Verify 함수를 호출
	//true 서명이 아직 유효하다는 것을 확인
	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)
	fmt.Println(ok)
}

/*
//* privateKey publickey signiture 생성 방법
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

	//따라서 먼저 생성한 비공개키를 숫자로 병경
	//비공개키를 특정 포맷으로 변경해서 저장할 수 있게 해주는 패키지 존재
	//앞으로 비공개키는 우리 파일 시스템에 저장된 file로 부터 가져올 것이다
	keyAsBytes , err := x509.MarshalECPrivateKey(privateKey)
	fmt.Printf("%x\n\n",keyAsBytes)

	//message := "I love you"
	//hashMessage := utils.Hash(message)
	//fmt.Println(hashMessage)
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
	//fmt.Printf("R:%d\nS:%d", r, s)
	//서명은 R 과 S 값으로 나눠져 있다.
	//r 과 s 는 big.Int 형식인데  big.Int 를 bytes 로 바꿔주는 메서드 있음
	//fmt.Println(r.Bytes(),s.Bytes())
	//r 과 s 통합 방법? slice 합치기
	//append 함수는 slcide 에 요소를 추가하는데
	// 두번째 매개변수에 element 를 추가하는데
	//... 은 elements 를 slice 로 부터 꺼내는 방법
	signature := append(r.Bytes(),s.Bytes()...)
	//len(r.Bytes()) => 32
	//signature 은 32byres 로 된 두 개의 slcide 이다
	fmt.Printf("%x\n",signature)

	//함수를 실행할때마다 signature갑이 바뀐다.
	//그 이유는 우리가 이 Start 를 실행할 때마다 새로운 비공개키를 생성하기 때문
	//증명?검증? 이 메세지에 한 서명이 맞는지 틀린지 확인
	//ok := ecdsa.Verify(&privateKey.PublicKey,hashAsByte,r,s)
	//fmt.Println(ok)
}*/
