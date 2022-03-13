package utils

import (
	"bytes"
	"encoding/gob"
	"log"
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
