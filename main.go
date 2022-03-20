package main

import "github.com/gwiyeomgo/nomadcoin/wallet"

/*func main() {
	//port 2개 실행시 go 사용
	rest.Start(3000)
	//go explorer.Start(4000)
}*/

/*func main() {
	defer db.Close()
	//CLI 는 유저에게 flag 를 입력하도록 요청
	//blockchain.Blockchain()
	cli.Start()
	//		blockchain.Blockchain().AddBlock("First")
	//		blockchain.Blockchain().AddBlock("Second")
	//		blockchain.Blockchain().AddBlock("Third")
}
*/
//과제?
//두 가지 mode 를 동시에 구동시키는 command나 flag 를 추가하기?
//두 가지 mode,(html,rest) 실행시키는 하나의 command 를 만들자

func main() {
	wallet.Start()
}
