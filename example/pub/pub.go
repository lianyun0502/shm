package main

import (
	"fmt"
	"time"

	"github.com/lianyun0502/shm"
)

func main() {
	publisher := shm.NewPublisher(66, 1024*1024)
	defer publisher.Close()


	for i:=0; i<100000; i++ {
		s := fmt.Sprintf("hello %d", i)
		publisher.Write([]byte(s))
		fmt.Println("write : ", s)
		time.Sleep(time.Millisecond*10)
	}

}