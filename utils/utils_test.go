package utils

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

//go test ./...  모든 디렉토리의 모든 파일에 대해 테스트 ,-v  모든 과정이 출력되는 옵션 => go test ./... -v
//go test -v -coverprofile cover.out ./... => 모든 파일 전체 커러비지 알 수 있음 =>cover.out 파일
//go tool cover -html=cover.out
func TestHash(t *testing.T) {
	hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
	s := struct {
		Test string
	}{Test: "test"}
	//콘솔에 찍힘
	//t.Log(x)
	t.Run("Hash is always smae", func(t *testing.T) {
		x := Hash(s)
		if x != hash {
			t.Errorf("Expected %s, got %s", hash, x)
		}
	})
	t.Run("Hash is hex encoded", func(t *testing.T) {
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Errorf("Hash should be hex encoded")
		}
	})
}

func ExampleHash() {
	s := struct {
		Test string
	}{Test: "test"}
	x := Hash(s)
	fmt.Println(x)
	// Output e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746

}

func TestToBytes(t *testing.T) {
	s := "test"
	b := ToBytes(s)
	//reflect 는 타입을 체크한다
	k := reflect.TypeOf(b).Kind()
	//k가 sclide 인지 확인
	if k != reflect.Slice {
		t.Errorf("ToBytes should return a slice of bytes got %s", k)
	}
}
