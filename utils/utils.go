package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	//golang 의 value 를 byte encode 나 decode 하는 패키지 gob
	// encoder 를 만들고 block 을 encode 한다음
	// 그 결과를 blockBuffer 에 저장
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	HandleErr(encoder.Encode(i))
	return blockBuffer.Bytes()
}

//bytes decode
//db에서 찾은 byte를 텅빈 블록체인의 memory address 에 decode
func FromBytes(i interface{}, data []byte) {
	/*decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)*/
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}

func Hash(i interface{}) string {
	//inerface formate %v => string
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

//문자열을 분리하는 함수
func Splitter(s string, sep string, i int) string {
	//슬라이스 길이보다 큰 인덱스를 요청했나?
	//slice -1 의 길이가 우리가 원하는 인데스보다 작은지 확인
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}
