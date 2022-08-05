package blockchain

import (
	"github.com/gwiyeomgo/nomadcoin/utils"
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	//메모리품에 transaction 을 하나 임의로 생성해 준다
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	//블록이 만들어졌는지 확인
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a blocks")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}
		_, err := FindBlock("xx")
		if err == nil {
			t.Error("The block should not be found")
		}
	})
	t.Run("Block id found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 1,
				}
				return utils.ToBytes(b)
			},
		}
		b, _ := FindBlock("xx")
		if b.Height != 1 {
			t.Error("Block should be found")
		}
	})
}
