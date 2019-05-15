package main

import (
	"log"

	"github.com/chrsm/ghettowm/virtd"
	"github.com/chrsm/winapi/user"
)

const (
	hotkeyDecrement = 1
	hotkeyIncrement = 2
	hotkeyQuit      = 3
)

var (
	maxVirt = virtd.GetDesktopCount() - 1
	current = virtd.GetCurrentDesktopNumber()
)

func main() {
	log.Printf("Starting ghettowm..")

	// <- dec, -> inc
	user.RegisterHotkey(nil, hotkeyDecrement, user.ModAlt|user.ModNoRepeat, user.VirtKeyLeft)
	user.RegisterHotkey(nil, hotkeyIncrement, user.ModAlt|user.ModNoRepeat, user.VirtKeyRight)
	user.RegisterHotkey(nil, hotkeyQuit, user.ModAlt|user.ModShift|user.ModNoRepeat, 0x51)

	log.Printf("max desktop id: %d", maxVirt)

	for {
		msg, ok := user.GetMessage(nil, 0, 0)

		if !ok {
			log.Fatal("user.GetMessage failed")
		}

		if msg.Message != user.WmHotkey {
			continue
		}

		if msg.WParam == hotkeyQuit {
			log.Printf("got hotkeyQuit")
			break
		}

		// WParam = the ID of the hotkey that was pressed.
		if msg.WParam != hotkeyDecrement && msg.WParam != hotkeyIncrement {
			continue
		}

		direction := 1
		if msg.WParam == hotkeyDecrement {
			direction = -1
		}

		switchTo(bound(current+direction, maxVirt))
	}

	user.UnregisterHotkey(nil, hotkeyDecrement)
	user.UnregisterHotkey(nil, hotkeyIncrement)
	user.UnregisterHotkey(nil, hotkeyQuit)
}

func WinMain(wproc uintptr) {
	main()
}

func switchTo(next int) {
	if current == next {
		log.Printf("tried to switch to current desktop(%d)", current)
		return
	}

	log.Printf("switching desktops: %d to %d", current, next)

	virtd.GoToDesktopNumber(next)

	current = next
}

func bound(i, max int) int {
	// loop backwards
	if i > max {
		return 0
	}

	// loop forwards
	if i < 0 {
		return max
	}

	return i
}
