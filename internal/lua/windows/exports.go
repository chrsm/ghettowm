package windows

import (
	"github.com/chrsm/winapi"
	"github.com/chrsm/winapi/user"
	lua "github.com/yuin/gopher-lua"
)

func Preload(ls *lua.LState) {
	ls.PreloadModule("windows", loader)
}

var exports = map[string]lua.LGFunction{
	"get_active_window":     getActiveWindow,
	"get_foreground_window": getForegroundWindow,
	"set_foreground_window": setForegroundWindow,
	"set_focus":             setFocus,
}

func loader(ls *lua.LState) int {
	mod := ls.SetFuncs(ls.NewTable(), exports)
	ls.Push(mod)
	return 1
}

func getActiveWindow(ls *lua.LState) int {
	hwnd := uintptr(user.GetActiveWindow())

	ls.Push(lua.LNumber(hwnd))
	return 1
}

func getForegroundWindow(ls *lua.LState) int {
	hwnd := uintptr(user.GetForegroundWindow())

	ls.Push(lua.LNumber(hwnd))
	return 1
}

func setForegroundWindow(ls *lua.LState) int {
	hwnd := winapi.HWND(uintptr(ls.CheckInt(1)))

	ok := user.SetForegroundWindow(hwnd)
	ls.Push(lua.LBool(ok))

	return 1
}

func setFocus(ls *lua.LState) int {
	hwnd := winapi.HWND(uintptr(ls.CheckInt(1)))

	user.SetFocus(hwnd)

	return 0
}
