package main

import (
	"log"
	"runtime"

	"github.com/chrsm/ghettowm"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting ghettowm..")

	gwm := ghettowm.New()
	gwm.Run()
}

func WinMain(wproc uintptr) {
	main()
}
