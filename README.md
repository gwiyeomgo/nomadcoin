# golang basic

# Variables in Go
```go
var name string = "test"
// type 이 compiler 에게 name 은 항상 string 이라고 말해준다
name1 :="test"
//compiler 가 자동으로 name의 type이 string 이라는 걸 안다.
const name string =""

```

# What is type int in Golang?
int         :signed integer (정수)
,int 뒤의 숫자는 크기를 나타냄
ex> unsigned 8-bit integers (0 to 255)
unit        :unsigned integer(= 음의 정수 아닌 양의 정수만 해당)


[출처](https://go.dev/tour/basics/11)

# Functions
argument
```go
package main
import (
	"fmt"
)
func plus (array ...int)  int{
	// a는 int의 array
	results := 0
	for _,item :=range array {
		results += item
	}
	return results
}
func main()  {
	result := plus(1,1,1,1,1,2)
	fmt.Println(result)
	name :="gwiyeom go~~!!#!#!"
	for _,item :=range name{
		//fmt.Println(item) //byte 타입으로 출력된다.
		fmt.Println(string(item))
	}
}

````

```
multiple return 



# 'text/template' HTML templage 가 구현하는 package
https://pkg.go.dev/text/template

main에서 
pages,partials 폴더에 만든 template을
일일이 load 하지 않고
home 에 대한 요청이 있을 때 마다
home function 이 호출 될 때마다
home template 을 parsing 하고
main 에 load

main 코드에  variable 을 추가했는데
해당 variable은 전체 template을 관리한다

