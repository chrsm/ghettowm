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
		desktopCount:   virtd.GetDesktopCount(),
		currentDesktop: virtd.GetCurrentDesktopNumber(),
		keybinds: &keybinds{
			set: make(map[int]*keybind),
		},
		windows: make(map[int]*desktop),
	}

	for i := 0; i < gwm.desktopCount; i++ {
		log.Printf("init-ing dn(%d)", i)
		gwm.windows[i] = &desktop{}
	}

	gwm.run()
}

func WinMain(wproc uintptr) {
	main()
}
