package main

import (
	"fmt"
	"github.com/samuel/go-thrift"
	"net"
	"test"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1463")
	if err != nil {
		panic(err)
	}

	client := thrift.NewClient(thrift.NewFramedReadWriteCloser(conn, 0), thrift.NewBinaryProtocol(true, false))
	tst := test.UserStorageClient{client}
	res, err := tst.Store(&test.UserProfile{1, "john doe", "amazing"})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Success? Server response: %+v\n", res)
}
