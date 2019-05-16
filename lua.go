package main

import (
	"github.com/chrsm/ghettowm/internal/lua/util"
	"github.com/chrsm/ghettowm/internal/lua/windows"

	"github.com/BixData/gluabit32"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func newLuaVM(gwm *ghettoWM) *lua.LState {
	opts := lua.Options{
		IncludeGoStackTrace: true,
	}

	ls := lua.NewState(opts)

	gluabit32.Preload(ls)

	windows.Preload(ls)
	util.Preload(ls)

	ls.SetGlobal("ghettowm", luar.New(ls, gwm))

	return ls
}
