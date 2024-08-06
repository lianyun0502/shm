package main

import (
	"fmt"
	"github.com/lianyun0502/shm"
)

func main() {
	subscriber := shm.NewSubscriber(66, 1024*1024)
	defer subscriber.Close()

	subscriber.Handle = func(data []byte) {
		fmt.Println("read: ", string(data))
		if string(data) == "EOF" {
			return
		}
	}
	subscriber.ReadLoop()
}
