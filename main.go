package main

import (
	_ "github.com/edwardOWO/goexample/learn"
	"github.com/edwardOWO/goexample/member"
	_ "github.com/edwardOWO/goexample/msg"
)

func main() {

	//msg2 := msg.NewMessage()
	//test := *(msg2)
	//fmt.Printf(test.Data)

	//msg2.Data = "123"

	//learn.GetReflect(*msg2)

	/*
		fmt.Println(t)

		msg2.SentMessage("test")
		msg2.SentMessage("2")

		msg3 := msg.NewMessage()
		msg3.SentMessage("2")
	*/

	//var test int = 10

	//learn.PrintValue(test)

	//go run main.go --members="192.168.1.104:40001,192.168.1.104:40002" --port=4002 --p=40002

	//172.17.0.3
	member.Main()

}
