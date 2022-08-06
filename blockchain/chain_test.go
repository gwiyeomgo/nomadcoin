package blockchain

import (
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/utils"
	"reflect"
	"sync"
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
	//once.Do 를 실행시키기 위해서 어떻게 하지?
	t.Run("Should restore blockchain", func(t *testing.T) {
		once = *new(sync.Once)
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{Height: 2, NewestHash: "xxx", CurrentDifficulty: 1}
				return utils.ToBytes(bc)
			},
		}
		bc := Blockchain()
		if bc.Height != 2 {
			t.Errorf("Blockchain() should restore a blockchain with a height of %d,got %d", 2, bc.Height)
		}
	})
}

func TestBlocks(t *testing.T) {
	//PreHash 를 return 하지 않을 때
	fakeBlocks := 0
	blocks := []*Block{
		{
			PreHash: "x",
		},
		{
			PreHash: "",
		},
	}
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() {
				fakeBlocks++
			}()
			/*var b *Block
			//PreHash 가 포함해서 return하면 Genesis 에
			if fakeBlocks == 0 {
				b = &Block{
					Height: 1,
					PreHash: "x",
				}
			}
			if fakeBlocks == 1 {
				//PreHash 가 있따면 Genesis 에 도착착
				b = &Block{
					Height: 1,
				}
			}
			fakeBlocks++*/
			return utils.ToBytes(blocks[fakeBlocks])
		},
	}
	bc := &blockchain{}
	blocksResult := Blocks(bc)
	if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
		t.Error("Blocks() should return a slice of blocks")
	}

}

func TestFindTx(t *testing.T) {
	//transcation 이 없는 block
	t.Run("Tx not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height:      2,
					Transaction: []*Tx{},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "x"}, "test")
		if tx != nil {
			t.Error("Tx should be not found")
		}
	})
	t.Run("Tx should be found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 2,
					Transaction: []*Tx{
						{ID: "test"},
					},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "x"}, "test")
		if tx == nil {
			t.Error("Tx should be not found")
		}
	})
}

//blockchain 의 height 가 0 이면 defaultDifficulty returnxxx
//다른 height 이면 CurrentDifficultly 를 return

func TestGetDifficulty(t *testing.T) {
	blocks := []*Block{
		{PreHash: "xxx"},
		{PreHash: "xxx"},
		{PreHash: "xxx"},
		{PreHash: "xxx"},
		{PreHash: ""},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() {
				fakeBlock++
			}()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}
	type test struct {
		height int
		want   int
	}
	tests := []test{
		{height: 0, want: defaultDifficulty},
		{height: 2, want: defaultDifficulty},
		{height: 5, want: 3},
	}
	for _, tc := range tests {
		bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDifficulty}
		b := getDifficulty(bc)
		if b != tc.want {
			t.Errorf("getDifficulty() should return %d got %d", b, tc.want)
		}
	}
}

func TestAddPeerBlock(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}
	m.Txs["test"] = &Tx{}
	b := &Block{
		Difficulty: 2,
		Hash:       "test",
		Transaction: []*Tx{
			{ID: "test"},
		},
	}
	bc.AddPeerBlock(b)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock should mutate the blockchain")
	}

}

func TestReplace(t *testing.T) {
	blocks := []*Block{
		{Difficulty: 2, Hash: "test"},
		{Difficulty: 2, Hash: "test"},
	}
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}
	bc.Replace(blocks)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("Replace() should mutate the blockchain")
	}

}
