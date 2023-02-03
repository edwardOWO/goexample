package msg

import "fmt"

var count int64 = 0

type msg struct {
	id      int64
	content string
}

func NewMessage() *msg {
	m := msg{}
	count += 1
	m.id = count
	m.content = "test"
	return &m
}

func NewMessage2() {
	fmt.Print("test")
}

func (msg *msg) SentMessage(data string) {
	fmt.Print(data)
	msg.content = data
}
