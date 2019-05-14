package main

import (
	"log"

	"github.com/chrsm/winapi"
	"github.com/chrsm/winapi/user"
)

const (
	hotkeyDecrement = 1
	hotkeyIncrement = 2
	hotkeyQuit      = 3

	maxVirt = 3
)

var (
	desktops map[int]map[winapi.HWND]struct{} = make(map[int]map[winapi.HWND]struct{})

	current = 0
)

func main() {
	log.Printf("Starting ghettowm..")

	// <- dec, -> inc
	user.RegisterHotkey(nil, hotkeyDecrement, user.ModAlt|user.ModNoRepeat, user.VirtKeyLeft)
	user.RegisterHotkey(nil, hotkeyIncrement, user.ModAlt|user.ModNoRepeat, user.VirtKeyRight)
	user.RegisterHotkey(nil, hotkeyQuit, user.ModAlt|user.ModShift|user.ModNoRepeat, 0x51)

	// windows per virtual desktop
	for i := 0; i < maxVirt; i++ {
		desktops[i] = make(map[winapi.HWND]struct{})
	}

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

	update()

	// hide old
	for hwnd := range desktops[current] {
		user.ShowWindow(hwnd, user.SwHide)
	}
	// show new
	for hwnd := range desktops[next] {
		user.ShowWindow(hwnd, user.SwShow)
	}

	current = next
}

func update() {
	for i := range desktops {
		vd := desktops[i]

		// remove dead windows
		for hwnd := range vd {
			handle, _ := user.GetWindowThreadProcessId(hwnd)
			if handle != 0 {
				continue
			}

			// remove this from the list
			delete(vd, hwnd)
		}
	}

	// and update the current desktop's state
	vd := desktops[current]
	for hwnd := range vd {
		if user.IsWindowVisible(hwnd) {
			continue
		}

		delete(vd, hwnd)
	}

	// finally, look for new ones
	user.EnumWindows(enumer, 0)
}

func enumer(window winapi.HWND, v winapi.LPARAM) uintptr {
	inf, err := user.GetWindowInfo(window)
	if err != nil {
		return 1
	}

	if inf.DwStyle&winapi.DWORD(user.WsVisible) != winapi.DWORD(user.WsVisible) || inf.DwExStyle&winapi.DWORD(user.WsExToolWindow) == winapi.DWORD(user.WsExToolWindow) {
		return 1
	}

	for i := range desktops {
		vd := desktops[i]

		for hwnd := range vd {
			if hwnd == window {
				return 1
			}
		}
	}

	desktops[current][window] = struct{}{}
	return 1
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
