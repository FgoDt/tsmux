package main

import (
	"fmt"

	"github.com/FgoDt/tsmux/decoder"
)

type testio string

func main() {
	d := &decoder.Decoder{}
	d.TSIO = new(testio)
	decoder.Decode(d)
	fmt.Println("run test")
}

func (i *testio) Read(data []byte, num int) error {
	data[0] = 0x47
	data[1] = 0x46
	data[2] = 0
	data[3] = 0x17
	return nil
}

func (i testio) Write(data []byte, num int) error {

	return nil
}
