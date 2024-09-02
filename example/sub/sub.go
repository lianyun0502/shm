package main

import (
	"fmt"
	"time"
	"github.com/lianyun0502/shm"
)

func main() {
	subscriber := shm.NewSubscriber(333, 1024*1024)
	defer subscriber.Close()

	handle := func(data []byte) {
		fmt.Println(fmt.Sprintf("read: %s time %d", string(data), time.Now().UnixNano()))
		if string(data) == "EOF" {
			return
		}
	}
	subscriber.Handle = handle
	subscriber.ReadLoop()
}
