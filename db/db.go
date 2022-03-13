package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gwiyeomgo/nomadcoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

//(1)Singleton 패턴을 사용해서 export 되지 않는 벼누를 마들고
var db *bolt.DB

//(2) export 되는 DB 함수를 만들었다
//db initialize function
func DB() *bolt.DB {
	if db == nil {
		//DB로 연결 open, db와 상호작용 할 수 있도록
		//1.make bucket = sql 에 table 과 같은 존재
		dbPointer, err := bolt.Open(dbName, 0600, nil)
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

func SaveBlock(hash string, data []byte) {
	fmt.Printf("Saving Block %s\nData: %b\n", hash, data)
	// (2)blockchain DB를 설정,생성
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		//key/value 추가
		err := bucket.Put([]byte("checkpoint"), data)
		return err
	})
	utils.HandleErr(err)
}
