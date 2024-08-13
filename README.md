# shm_agent

[![Go Report Card](https://goreportcard.com/badge/github.com/lianyun0502/shm)](https://goreportcard.com/report/github.com/lianyun0502/shm)
[![GoDoc](https://godoc.org/github.com/lianyun0502/shm?status.svg)](https://godoc.org/github.com/lianyun0502/shm)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)]

shm is a Pub/Sub pacakge that uses shared memory to communicate between processes.

## Index
- [Introduction](#Introduction)
- [Installation](#Installation)
- [Example](#Example)

## Installation

there are two ways to install and use the package, one is to clone the repository and refer to local path and the other is to use the `go get` command set to `go.mod`.

### Git clone the repository

1. first clone the repository into your project directory

    ```bash
    git clone https://github.com/lianyun0502/exchange_conn.git
    ```

    Your directory structure should look like this:

    ```bash
    your_project/
    ├── exchange_conn/
    ├── main.go
    └── go.mod
    ```
    
2. replace the import refernce with the path of the repository in your project.

    ```bash
    go mod edit -replace=github.com/lianyun0502/exchange_conn=../exchange_conn
    ```

3. import the package in your project

    ```Go
    import (
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )
    ```

### Install the package use `go get`

1. use the `go get` command to install the package

    ```bash
    go get github.com/lianyun0502/exchange_conn
    ```
2. import the package in your project

    ```Go
    import (
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )
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
 