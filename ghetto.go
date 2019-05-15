package main

import (
	"log"
	"os"

	"github.com/chrsm/ghettowm/virtd"
	"github.com/chrsm/winapi/user"
)

type ghettoWM struct {
	desktopCount   int
	currentDesktop int

	keybinds *keybinds
}

func (gwm *ghettoWM) RegisterHotkey(modifiers user.ModKey, vkey user.VirtualKey, cb func(int)) bool {
	nextID := gwm.keybinds.len + 1

	gwm.keybinds.set[nextID] = &keybind{
		id: nextID,
		cb: cb,
	}

	gwm.keybinds.len++

	return user.RegisterHotkey(nil, nextID, modifiers, vkey)
}

func (gwm *ghettoWM) SwitchDesktop(dst int) bool {
	if dst == gwm.currentDesktop {
		return false
	}

	virtd.GoToDesktopNumber(dst)
	gwm.currentDesktop = dst

	return true
}

func (gwm *ghettoWM) SwitchDesktopPrev() bool {
	return gwm.SwitchDesktop(bound(gwm.currentDesktop-1, gwm.desktopCount))
}

func (gwm *ghettoWM) SwitchDesktopNext() bool {
	return gwm.SwitchDesktop(bound(gwm.currentDesktop+1, gwm.desktopCount))
}

func (gwm *ghettoWM) Quit() {
	for k := range gwm.keybinds.set {
		kb := gwm.keybinds.set[k]

		user.UnregisterHotkey(nil, kb.id)
	}

	log.Println("Quitting!")
	os.Exit(0)
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
