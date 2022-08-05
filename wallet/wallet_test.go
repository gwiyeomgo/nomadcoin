package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey     string = "307702010104201faf5e112e0cdab799d95773e4b685fe12672e3487d1f51923989b92a36c6b61a00a06082a8648ce3d030107a14403420004655b79f1927d71f0b42f33b89a0dad3234b48e6403216456c24a9f26accaa5346082d833a20a9253e4fa79d98396b9250527b53dd42d89d3106d42a8988d855b"
	testPayload string = "00bf5fe3c6ab2ecf867f1387e6c92ac0b8dcb26da96ea57a2a7aa0d508f42b57"
	testSig     string = "fb495a87d0a411fbd3d23a3d722fc1b364d6ad2d959a700559865f40ec4081c78a2a648ae259946ade6b97d2d817f5a43099280446cd488c13a0b2fabc424f82"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	//testkey 를 decode 하고
	b, _ := hex.DecodeString(testKey)
	//private key 로 복원한다
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestSign(t *testing.T) {
	//뭔가에 서명해야 하고 그 서명을 저장한다
	//서명에 대한 검증을 할 수 있다
	s := Sign(testPayload, makeTestWallet())
	//t.Log(s) //testSig 값
	//서명은 hash 와 달리 램덤으로 만들어진다 -> 일치 테스트 불가
	//그래서 hex.DecodeString
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string,got %s", s)
	}
}
func TestVerify(t *testing.T) {
	/*	//test key 를 얻기위한 코드
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
		b,_:= x509.MarshalECPrivateKey(privateKey)
		//t.Logf()를 사용해서 hexadecimal 로 key 를 볼 수 있다.
		t.Logf("%x",b)*/
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"00bf5fe3c6ab2ecf867f1387e6c92ac0b8dcb26da96ea57a2a7aa0d508f42b257", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts should return error when payload is not hex")
	}
}

type fakeLayer struct {
	fakeHashWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHashWalletFile()
}
func (f fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

//wqllet 의 bytes 를 return 해준다
func (f fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func TestWallet(t *testing.T) {
	//wallet 을 테스트함
	//그런데 files 의 의미를 변경해준다
	//Layer 를 가짜 Layer로 변경

	//1. wallet 파일이 있는지 확인,있을 때는 이것들이 작동하는지, wallet 을 return 하는지 확인
	//2. wallet 파일이 없다면, privateKey 생성해 wallet 을 생성하는지 확인
	t.Run("New Wallet is createde", func(t *testing.T) {
		files = fakeLayer{
			fakeHashWalletFile: func() bool { return false },
		}
		//wallet 이 만들어졌을 때
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return  a new wallet instance")
		}
	})
	t.Run(" Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHashWalletFile: func() bool { return true },
		}
		//wallet 이 만들어졌을 때
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return  a new wallet instance")
		}

	})
}
