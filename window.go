package main

import (
	"github.com/chrsm/ghettowm/virtd"
	"github.com/chrsm/winapi"
)

const (
	gwmWindowFloat = 0x01
	gwmWindowTile  = 0x02
)

type window struct {
	hwnd winapi.HWND

	// our personal state
	floating bool

	next, prev *window
}

func (w *window) DesktopNumber() int {
	return virtd.GetWindowDesktopNumber(w.hwnd)
}

func (w *window) MoveTo(dst int) bool {
	return virtd.MoveWindowToDesktopNumber(w.hwnd, dst)
}

func (w *window) IsPinned() bool {
	return virtd.IsPinnedWindow(w.hwnd) // || virtd.IsPinnedApp(w.hwnd)
}

func (w *window) Pin() {
	virtd.PinWindow(w.hwnd)
}
