package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/chrsm/ghettowm/virtd"
	"github.com/chrsm/winapi"
	"github.com/chrsm/winapi/kernel"
	"github.com/chrsm/winapi/user"
)

type window struct {
	id int

	hwnd winapi.HWND
}

type ghettoWM struct {
	desktopCount   int
	currentDesktop int

	// per-desktop window count.
	windows map[int][]*window

	keybinds *keybinds

	hwnd winapi.HWND
}

func (gwm *ghettoWM) RegisterHotkey(modifiers user.ModKey, vkey user.VirtualKey, cb func(int)) bool {
	gwm.keybinds.len++

	gwm.keybinds.set[gwm.keybinds.len] = &keybind{
		id: gwm.keybinds.len,
		cb: cb,
	}

	return user.RegisterHotkey(gwm.hwnd, gwm.keybinds.len, modifiers, vkey)
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

		user.UnregisterHotkey(gwm.hwnd, kb.id)
	}

	log.Println("Quitting!")
	os.Exit(0)
}

func (gwm *ghettoWM) run() {
	const HWND_MESSAGE = winapi.HWND(^uintptr(2))
	cname := "ghettowm"

	self := kernel.GetModuleHandle("")

	wclass := user.WindowClass{
		WndProc:     syscall.NewCallback(gwm.wndproc),
		HInstance:   winapi.HINSTANCE(self),
		ClassName:   syscall.StringToUTF16Ptr(cname),
		HBackground: 6,
	}
	wclass.CbSize = winapi.DWORD(unsafe.Sizeof(wclass))

	clsid := user.RegisterClass(&wclass)
	hwnd := user.CreateWindow(
		cname,
		cname,
		0,
		0,
		0,
		0,
		0,
		HWND_MESSAGE,
		0,
		winapi.HMODULE(self),
		nil,
	)

	gwm.hwnd = hwnd

	msgnum := user.RegisterWindowMessage("SHELLHOOK")
	user.RegisterShellHookWindow(hwnd)

	vm := newLuaVM(gwm)
	defer vm.Close()

	// Run user configuration through lua, because I don't feel that
	// writing a conf language for this makes sense, and opens up more
	// customization options in the future.
	if err := vm.DoFile("ghetto.lua"); err != nil {
		panic(err)
	}

	log.Println("configuration succeeded")

	msg := &user.Message{}
	for {
		if ok := user.GetMessage(msg, 0, 0, 0); !ok {
			log.Fatal("/shrug")
		}

		user.TranslateMessage(msg)
		user.DispatchMessage(msg)
	}

	_ = clsid
	_ = msgnum
}

func (gwm *ghettoWM) wndproc(hwnd winapi.HWND, msg uint32, wparam uintptr, lparam uintptr) uintptr {
	log.Printf("hwnd(%X); msg(%X, %d); wp(%X, %d); lp(%X, %d)", hwnd, msg, msg, wparam, wparam, lparam, lparam)
	switch msg {
	case user.WmHotkey:
		log.Printf("WM_HOTKEY")

		if kb, ok := gwm.keybinds.set[int(wparam)]; ok {
			kb.cb(int(wparam))
		} else {
			return user.DefWindowProc(hwnd, msg, winapi.WPARAM(wparam), winapi.LPARAM(lparam))
		}
	default:
		return user.DefWindowProc(hwnd, msg, winapi.WPARAM(wparam), winapi.LPARAM(lparam))
	}

	return 0
}
