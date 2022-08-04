package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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

func TestSplitter(t *testing.T) {
	//테이블 테스트란
	//테스트 function 맨 위에 많은 테스트를 정의해놓고 그 모든 테스트를 for loop 로 매우 빠르게 수행한다
	//t.Run 을 반복할 필요는 없다
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	//golang 의 test 는 사람들이 input,output 을 정의하고 어떤게 나오는지,또 어떤걸 받아야하는지 확인하는걸 볼 수 있다
	tests := []test{
		{input: "0:6:0", sep: ":", index: 1, output: "6"},
		{input: "0:6:0", sep: ":", index: 10, output: ""},
		{input: "0:6:0", sep: "/", index: 10, output: "0:6:0"},
	}
	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if got != tc.output {
			t.Errorf("Expected %s and got %s", tc.output, got)
		}
	}
}

func TestHandleErr(t *testing.T) {
	oldLogFn := logFn
	//테스트가 끝나면 백업한 원래 logFn 함수로 바꾼다
	defer func() {
		logFn = oldLogFn
	}()
	called := false
	logFn = func(v ...interface{}) {
		called = true
	}
	err := errors.New("test")
	HandleErr(err)
	if !called {
		t.Error("HandleError should call fn")
	}
}

func TestFromBytes(t *testing.T) {
	type testStruct struct {
		Test string
	}
	var restored testStruct
	ts := testStruct{"test"}
	b := ToBytes(ts)
	//텅빈 restored 와 []byte 를 넘겨서
	FromBytes(&restored, b)
	if !reflect.DeepEqual(ts, restored) {
		t.Error("FromBytes() should restore struct.")
	}
}

func TestToJSON(t *testing.T) {
	type TestStruct struct{ Test string }
	s := TestStruct{"test"}
	//json 을 bytes 로 이코딩
	b := ToJSON(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		//reflect.Slice 를 기대했지만 reflect.TypeOf(got).Kind() 이 왔다
		t.Errorf("Expected %v and got %v", reflect.Slice, k)
	}
	//인코됭된 bytes 를 json 으로 복구
	var restored TestStruct
	json.Unmarshal(b, &restored)
	if !reflect.DeepEqual(s, restored) {
		//JSON을 올바르게 encode 해야 한다
		t.Errorf("ToJSON() should encode to JSON correctly")
	}

}
