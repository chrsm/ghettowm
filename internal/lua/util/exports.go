package util

import (
	"github.com/chrsm/winapi/user"
	lua "github.com/yuin/gopher-lua"
)

func Preload(ls *lua.LState) {
	ls.PreloadModule("ghetto_util", loader)
}

var exports = map[string]lua.LGFunction{
	"get_key": getKeyCode,
	"get_mod": getModCode,
}

func loader(ls *lua.LState) int {
	mod := ls.SetFuncs(ls.NewTable(), exports)
	ls.Push(mod)
	return 1
}

func getKeyCode(ls *lua.LState) int {
	name := ls.CheckString(1)
	kc := user.GetVirtualKeyByName(name)
	ls.Push(lua.LNumber(kc))

	return 1
}

func getModCode(ls *lua.LState) int {
	name := ls.CheckString(1)
	kc := user.GetModKeyByName(name)
	ls.Push(lua.LNumber(kc))

	return 1
}
