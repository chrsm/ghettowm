package main

import (
	"log"
	"runtime"

	"github.com/chrsm/ghettowm/virtd"
)

func main() {
	runtime.LockOSThread()

	log.Println("Starting ghettowm..")

	gwm := &ghettoWM{
		desktopCount:   virtd.GetDesktopCount() - 1,
		currentDesktop: virtd.GetCurrentDesktopNumber(),
		keybinds: &keybinds{
			set: make(map[int]*keybind),
		},
	}

	gwm.run()
}

func WinMain(wproc uintptr) {
	main()
}
