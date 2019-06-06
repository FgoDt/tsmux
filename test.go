package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/FgoDt/tsmux/decoder"
)

type testio string

var f *os.File

func main() {
	var e error
	f, e = os.OpenFile("bipbop.ts", os.O_RDONLY, 0666)
	if e != nil {
		fmt.Print(e.Error())
	}

	d := &decoder.Decoder{}
	d.TSIO = new(testio)
	d.Init()
	d.Run()
	fmt.Println("run test")

}

func (i *testio) Read(data []byte, num uint64) error {
	len, e := f.Read(data)
	if uint64(len) < num || e != nil {
		return errors.New("no more data")
	}
	return nil
}

func (i testio) Write(data []byte, num uint64) error {

	return nil
}
