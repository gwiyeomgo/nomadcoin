package main

import "github.com/gwiyeomgo/nomadcoin/cli"

/*func main() {
	//port 2개 실행시 go 사용
	rest.Start(3000)
	//go explorer.Start(4000)
}*/

func main() {
	//CLI 는 유저에게 flag 를 입력하도록 요청
	cli.Start()
}

//과제?
//두 가지 mode 를 동시에 구동시키는 command나 flag 를 추가하기?
//두 가지 mode,(html,rest) 실행시키는 하나의 command 를 만들자
