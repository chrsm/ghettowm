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

const (
	windowNext = 1
	windowPrev = 2
)

type ghettoWM struct {
	desktopCount   int
	currentDesktop int

	// per-desktop windows
	windows map[int]*desktop

	keybinds *keybinds

	hwnd  winapi.HWND
	shmsg uint

	rwnd winapi.HWND
}

func (gwm *ghettoWM) RegisterHotkey(modifiers user.ModKey, vkey user.VirtualKey, cb func(int)) bool {
	gwm.keybinds.len++

	gwm.keybinds.set[gwm.keybinds.len] = &keybind{
		id: gwm.keybinds.len,
		cb: cb,
	}

	return user.RegisterHotkey(gwm.hwnd, gwm.keybinds.len, modifiers, vkey)
}

func (gwm *ghettoWM) switchWindow(direction int) {
	curdeskn := virtd.GetCurrentDesktopNumber()
	curdesk, ok := gwm.windows[curdeskn]
	if !ok {
		log.Printf("Tried to switch windows on a desktop that doesn't exist internally (%d)", curdeskn)
	}

	cur := curdesk.find(user.GetForegroundWindow())
	if cur == nil {
		log.Printf("could not find active window(%X)", int(user.GetForegroundWindow()))
		return
	}

	var dst *window

	switch direction {
	case windowNext:
		dst = cur.next

		if dst == nil {
			// we've reached the end, so we go back to the first window.
			dst = curdesk.head
		}
	case windowPrev:
		dst = cur.prev
		if dst == nil {
			// we've reached the beginning, so we go back to the last window.
			dst = curdesk.tail
		}
	default:
		panic("unknown direction")
	}

	if dst == nil {
		log.Printf("no window to switch to, unfortunately..")
		return
	}

	user.SetForegroundWindow(dst.hwnd)
}

func (gwm *ghettoWM) NextWindow() {
	gwm.switchWindow(windowNext)
}

func (gwm *ghettoWM) PrevWindow() {
	gwm.switchWindow(windowPrev)
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

	gwm.rwnd = user.GetDesktopWindow()
	log.Printf("desktop window(%X, %d); owner(%X)", gwm.rwnd, gwm.rwnd, user.GetWindow(gwm.rwnd, user.GwOwner))
	self := kernel.GetModuleHandle("")

	wclass := user.WindowClass{
		WndProc:   syscall.NewCallback(gwm.wndproc),
		HInstance: winapi.HINSTANCE(self),
		ClassName: syscall.StringToUTF16Ptr(cname),
	}
	wclass.CbSize = winapi.DWORD(unsafe.Sizeof(wclass))

	user.RegisterClass(&wclass)
	gwm.hwnd = user.CreateWindow(
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

	// I'm not sure if this is always 0xC028..
	gwm.shmsg = user.RegisterWindowMessage("SHELLHOOK")
	user.RegisterShellHookWindow(gwm.hwnd)

	vm := newLuaVM(gwm)
	defer vm.Close()

	// Run user configuration through lua, because I don't feel that
	// writing a conf language for this makes sense, and opens up more
	// customization options in the future.
	if err := vm.DoFile("ghetto.lua"); err != nil {
		panic(err)
	}

	log.Println("configuration succeeded")

	user.EnumWindows(gwm.enumproc, 0)

	msg := &user.Message{}
	for {
		if ok := user.GetMessage(msg, 0, 0, 0); !ok {
			log.Fatal("/shrug")
		}

		user.TranslateMessage(msg)
		user.DispatchMessage(msg)
	}
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
		if msg == uint32(gwm.shmsg) && isUsefulHShellMsg(wparam) {
			// handle it
			return 0
		}

		return user.DefWindowProc(hwnd, msg, winapi.WPARAM(wparam), winapi.LPARAM(lparam))
	}

	return 0
}

func (gwm *ghettoWM) enumproc(hwnd winapi.HWND, _ winapi.LPARAM) uintptr {
	// check if the window is something we want to be responsible for

	buf := make([]uint16, 32)
	buf2 := make([]uint16, 32)
	// do something with it..
	title := user.GetWindowText(hwnd, buf)
	_, pid := user.GetWindowThreadProcessId(hwnd)
	cname := user.GetClassName(hwnd, buf2)

	log.Printf("hwnd(%X,%d); cname(%s); desktopn(%d); title(%s); pid(%d)", hwnd, hwnd, cname, virtd.GetWindowDesktopNumber(hwnd), title, pid)
	if usableWindow(hwnd) {
		gwm.windows[virtd.GetWindowDesktopNumber(hwnd)].push(hwnd)
	}

	return 1
}
