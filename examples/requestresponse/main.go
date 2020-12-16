package main

import (
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Hello struct{ Who string }

func Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		context.Respond("Hello " + msg.Who)
	}
}

func main() {
	rootContext := actor.EmptyRootContext
	props := actor.PropsFromFunc(Receive)
	pid := rootContext.Spawn(props)
	// 类似于RPC调用，可以设置超时时间
	result, _ := rootContext.RequestFuture(pid, Hello{Who: "Roger"}, 30*time.Second).Result() // await result

	fmt.Println(result)
	console.ReadLine()
}
