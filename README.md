# shm_agent

## Description
shm is a Pub/Sub pacakge that uses shared memory to communicate between processes.

## Index
- [Description](#description)

## Installation
* git clone the repository
```bash
git clone https://github.com/lianyun0502/shm.git
```
* go get the package
```bash
go get github.com/lianyun0502/shm
```

## Example

there are two examples in the example folder, one is the publisher and the other is the subscriber. set the handler in the subscriber to handle the data. the publisher writes data to the shared memory and the subscriber reads the data from the shared memory. The `skey` of the publisher and the subscriber should be the same. The `size of the shared memory` should be the same as well.

* publisher
 ```go
    package main

    import (
        "fmt"
        "time"

        "github.com/lianyun0502/shm"
    )

    func main() {
        // create a new publisher
        publisher := shm.NewPublisher(66, 1024*1024)
        // close the publisher when done
        defer publisher.Close()

        // write data to the shared memory
        for i:=0; i<100000; i++ {
            s := fmt.Sprintf("hello %d", i)
            publisher.Write([]byte(s))
            fmt.Println("write : ", s)
            time.Sleep(time.Millisecond*10)
        }

    }
```

* subscriber
 ```go
    package main

    import (
        "fmt"
        "github.com/lianyun0502/shm"
    )

    func main() {
        // create a new subscriber
        subscriber := shm.NewSubscriber(66, 1024*1024)
        // close the subscriber when done
        defer subscriber.Close()
        // set handler to handle the data
        subscriber.Handle = func(data []byte) {
            fmt.Println("read: ", string(data))
            if string(data) == "EOF" {
                return
            }
        }
        // start the subscriber
        subscriber.ReadLoop()
    }
```
 