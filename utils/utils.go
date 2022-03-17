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

//bytes decode
//db에서 찾은 byte를 텅빈 블록체인의 memory address 에 decode
func FromBytes(i interface{}, data []byte) {
	/*decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)*/
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}
