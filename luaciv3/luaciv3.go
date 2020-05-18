package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// NewState is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

// NewState creates and returns a lua state with LuaCiv3 functions
// TODO: memory limit?
func NewState() *lua.LState {
	L := lua.NewState()
	// Set memory limit of lua instance (just a safety measure)
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}
	_ = LuaCiv3(L)
	return L
}

// LuaCiv3 injects functions into a gopher-lua state
// TODO: Should I eliminate error return? lua.NewState() doesn't return error
func LuaCiv3(L *lua.LState) error {
	// test function
	L.SetGlobal("test", L.NewFunction(TestPassValues))

	// civ3 table
	civ3 := L.NewTable()
	if path, err := findWinCivInstall(); err == nil {
		L.RawSet(civ3, lua.LString("path"), lua.LString(path))
	} else {
		return err
	}
	L.SetGlobal("civ3", civ3)
	return nil
}

// TestPassValues is my getting familiar with calling Go from lua with values/params
//  see https://github.com/yuin/gopher-lua#calling-go-from-lua
func TestPassValues(L *lua.LState) int {
	lv := L.ToInt(1)
	lv2 := L.ToInt(2)
	L.Push(lua.LNumber(lv * lv2))
	return 1
}
