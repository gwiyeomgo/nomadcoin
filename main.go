package main

import (
	"fmt"
	"time"
)

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
/*
func main() {
	wallet.Wallet()
}*/
//# ------------- 고루틴

/*func countToTen()  {
	for i := range [10]int{} {
		fmt.Println(i)
		time.Sleep(1 * time.Second)

	}
}*/

//countToTen 밖에서 값을 얻고 싶다면?
// go 루틴 밖에서 값을 얻는 방법?
// * go 루틴 함수는 return 불가

func countToTen() {
	for i := range [10]int{} {
		fmt.Println(i)
		time.Sleep(1 * time.Second)

	}
}
func main() {
	// 아래의 경우 하나의 function 이 끝나기 전까지 다음 function은 실행되지 않는다
	//countToTen()
	//countToTen()
	//	쓴다면 순차적으로 0~9 까지 두번 출력됨
	go countToTen()
	go countToTen()
	// main function 이 고루틴을 기다려주지 않기떄문에
	// main 은 main 내부 함수가 끝날 때까지 끝나지 않지만
	//go 를 붙이면 go언어에게 기다리지 말라고 말해준다
	for {

	}
}
