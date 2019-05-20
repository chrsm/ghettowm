package main

import (
	"github.com/chrsm/ghettowm/virtd"
	"github.com/chrsm/winapi"
	"github.com/chrsm/winapi/user"
)

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

func isUsefulHShellMsg(wparam uintptr) bool {
	if wparam == user.HShellWindowActivated ||
		wparam == user.HShellWindowCreated ||
		wparam == user.HShellWindowDestroyed ||
		wparam == user.HShellWindowReplaced {
		return true
	}

	return false
}

func usableWindow(hwnd winapi.HWND) bool {
	// is it even visible? lol
	if !user.IsWindowVisible(hwnd) || !user.IsWindow(hwnd) {
		return false
	}

	// There are some "windows" that aren't actually on the desktop.
	// I don't know enough about Win32 to say _why_ this is.
	// For example, there's an invisible Calculator window. It seems
	// to become active after the app is launched.
	if virtd.GetWindowDesktopNumber(hwnd) == -1 {
		return false
	}

	// we need to remember to re-check the parent after this..
	parent := user.GetParent(hwnd)
	// src := user.GetWindow(hwnd, user.GwOwner)

	style, exstyle := user.GetWindowLong(hwnd, user.GwlStyle), user.GetWindowLong(hwnd, user.GwlExStyle)

	if style&user.WsDisabled == user.WsDisabled {
		return false
	}

	// ..figure out how to handle these
	if exstyle&user.WsExToolWindow == user.WsExToolWindow {
		return false
	}

	buf := make([]uint16, 10)
	if len(user.GetWindowText(hwnd, buf)) == 0 {
		return false
	}

	return parent == 0
}
