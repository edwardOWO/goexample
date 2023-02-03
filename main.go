package main

import (
	"github.com/edwardOWO/goexample/msg"
	_ "github.com/edwardOWO/goexample/msg"
)

func main() {

	msg2 := msg.NewMessage()
	msg2.SentMessage("test")
	msg2.SentMessage("2")

	msg3 := msg.NewMessage()
	msg3.SentMessage("2")

}
