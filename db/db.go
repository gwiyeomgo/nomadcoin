package db

import (
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/utils"
	bolt "go.etcd.io/bbolt"
	"os"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkPoint   = "checkpoint"
)

//(1)Singleton 패턴을 사용해서 export 되지 않는 벼누를 마들고
var db *bolt.DB

//현재 port 를 알 수 있는 함수
func getPort() string {
	/*	for i, a :=  range os.Args{
		//0
		//1 -mode=rest
		//2 -port=4000
		fmt.Println(i,a)
	}*/
	//fmt.Println(os.Args[2][6:])
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

//(2) export 되는 DB 함수를 만들었다
//db initialize function
func DB() *bolt.DB {
	if db == nil {

		//* open 이후에는 자료를 해방해 주고,
		//data 손상을 방지하기 위해 DB를 닫아줘야 한다.

		//DB로 연결 open, db와 상호작용 할 수 있도록
		//1.make bucket = sql 에 table 과 같은 존재
		dbPointer, err := bolt.Open(getPort(), 0600, nil)
		//path db name,mode 는 write/read 해야하는 permission
		db = dbPointer
		utils.HandleErr(err)
		//bucket 이 있는지 확인 -> data 저장 bucket,블록 저장 bucket
		err = db.Update(func(tx *bolt.Tx) error {
			//tx.CreatedBucketIfNotExists()
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}
func SaveBlock(hash string, data []byte) {
	//fmt.Printf("Saving Block %s\nData: %b\n", hash, data)
	// (2)blockchain DB를 설정,생성
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveCheckpoint(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		//key/value 추가
		err := bucket.Put([]byte(checkPoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkPoint))
		//bucket 은 err 반환 X
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	//DB에 blocksBucket 에서 특정 블록을 찾는다
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
