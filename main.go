package main

import (
	"log"

	"github.com/chrsm/ghettowm/virtd"
	"github.com/chrsm/winapi/user"
)

func main() {
	log.Println("Starting ghettowm..")

	gwm := &ghettoWM{
		desktopCount:   virtd.GetDesktopCount() - 1,
		currentDesktop: virtd.GetCurrentDesktopNumber(),
		keybinds: &keybinds{
			set: make(map[int]*keybind),
		},
	}

	vm := newLuaVM(gwm)
	defer vm.Close()

	// Run user configuration through lua, because I don't feel that
	// writing a conf language for this makes sense, and opens up more
	// customization options in the future.
	if err := vm.DoFile("ghetto.lua"); err != nil {
		panic(err)
	}

	log.Println("configuration succeeded")

	for {
		msg, ok := user.GetMessage(nil, 0, 0)
		if !ok {
			log.Fatal("failed to winapi/user.GetMessage")
		}

		if msg.Message != user.WmHotkey {
			continue
		}

		if kb, ok := gwm.keybinds.set[int(msg.WParam)]; ok {
			kb.cb(int(msg.WParam))
		}
	}
}

func WinMain(wproc uintptr) {
	main()
}
