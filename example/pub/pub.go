package main

import (
	"fmt"
	"time"

	"github.com/lianyun0502/shm"
	"github.com/sirupsen/logrus"
)

func main() {
	shm.Logger.SetLevel(logrus.DebugLevel)
	publisher := shm.NewPublisher(333, 1024*1024)
	publisher.Scheduler.Every(30).Second().Do(publisher.ResetMsgID)
	
	for i:=0; i<1000; i++ {
		s := fmt.Sprintf("hello %d, time %d", i, time.Now().UnixNano())
		publisher.Write([]byte(s))
		fmt.Println("write : ", s)
		time.Sleep(time.Millisecond*10)
	}
	publisher.Close()
}