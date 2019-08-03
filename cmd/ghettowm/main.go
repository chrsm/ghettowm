package main

import (
	"log"
	"runtime"

	"github.com/chrsm/ghettowm"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var runch = make(chan func())

	log.Println("Starting ghettowm..")

	gwm := ghettowm.New(runch)
	go gwm.Run()

	for f := range runch {
		f()
	}
}

func WinMain(wproc uintptr) {
	main()
}
