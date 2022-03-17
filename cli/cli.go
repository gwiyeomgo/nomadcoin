package cli

import (
	"flag"
	"fmt"
	"github.com/gwiyeomgo/nomadcoin/explorer"
	"github.com/gwiyeomgo/nomadcoin/rest"
	"runtime"
)

//CLI 적용
//command line과 code가 소통하도록 한다.
//ex) go run main.go rest -port 4000 명령하면
//rest API 가 port 4000 으로 시작된다.
func usage() {
	fmt.Printf("welcome to 노마드 코인\n")
	fmt.Printf("Please use the flollowing flags:\n\n")
	fmt.Printf("-port:	Set the PORT of the server\n")
	fmt.Printf("-mode:	Choose between 'html' and 'rest'\n\n")
	//프로그램을 종료시킴 => def 사용 후 종료되도록 수정
	//os.Exit(0)
	//GoExit 은 모든 함수를 제거하지만 그전에 defer 를 먼지 이행
	runtime.Goexit()
}
func Start() {
	//fmt.Println(os.Args)//string 문자열의 array
	//os.Args 의 길이가 2보다 작다면 입력값 x
	/*if len(os.Args) < 2 {
		usage()
	}*/
	port := flag.Int("port", 3000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")
	flag.Parse()
	//>go run main.go -mode rest -port 4000
	//fmt.Println(*port, *mode)
	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
	/*//[2:] 처럼 쓴다면 go 에서 array 의 요소 한 부분에서 끝까지 선택하는 방법
	fmt.Println(os.Args[2:])
	*/
	//FlagSet 은 go에게 어떤 command 가 어떤 flag 를 가질 것인지 알려준다
	//rest := flag.NewFlagSet("rest",flag.ExitOnError)
	//flag 가 많을때 FlagSet을 사용하는 것이 좋다
	//ex) go run main.go rest -port=4000 -mode https -v ....
	//port 값이 int 인지 확인
	/*portFlag := rest.Int("port",3000,"Sets the port of the server")
	switch os.Args[1] {
		case "explorer":
			fmt.Println("Start Explorer")
		case "rest":
			//go가 자동으로 Args에서 port 를 찾고
			rest.Parse(os.Args[2:])
		default:
			usage()
	}
	//fmt.Println(*portFlag)
	if rest.Parsed() {
		fmt.Printf("Start Server:%v",*portFlag)
	}*/

}
