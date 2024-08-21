package main

import (
	"fmt"
	"github.com/lianyun0502/shm"
	"github.com/sirupsen/logrus"
)



func main() {
	shm.Logger.SetLevel(logrus.DebugLevel)
	shm.Logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000", FullTimestamp: true})
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
