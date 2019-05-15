package main

import (
	"github.com/chrsm/winapi/user"

	"github.com/BixData/gluabit32"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func newLuaVM(gwm *ghettoWM) *lua.LState {
	opts := lua.Options{
		IncludeGoStackTrace: true,
	}

	state := lua.NewState(opts)
	gluabit32.Preload(state)

	state.SetGlobal("ghettowm", luar.New(state, gwm))
	state.SetGlobal("get_key_code", luar.New(state, user.GetVirtualKeyByName))
	state.SetGlobal("get_mod_code", luar.New(state, user.GetModKeyByName))

	return state
}
