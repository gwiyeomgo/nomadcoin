package blockchain

import (
	"fmt"
	"reflect"
	"testing"
)

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	fmt.Printf("hash:%v ", hash)
	return f.fakeFindBlock()
}

func (f fakeDB) SaveBlock(hash string, data []byte) {
	fmt.Printf("hash:%v and data:%x", hash, data)
}

func (f fakeDB) SaveChain(data []byte) {

}

func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}

func (f fakeDB) DeleteAllBlocks() {

}

func TestBlockchain(t *testing.T) {
	t.Run("Blockchain  not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}
		chain := Blockchain()
		if reflect.TypeOf(chain) != reflect.TypeOf(&blockchain{}) {
			t.Error("Blockchain not found")
		}
	})
}
