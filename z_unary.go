// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package lua

import (
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

func init() {
	typ := reflect.TypeOf((*func(String) error)(nil)).Elem()
	builtin[typ] = func(v any) lua.LGFunction {
		f := v.(func(String) error)
		return func(state *lua.LState) int {
			if err := f(String(state.CheckString(1))); err != nil {
				state.RaiseError(err.Error())
			}
			return 0
		}
	}
}

func init() {
	typ := reflect.TypeOf((*func() (String, error))(nil)).Elem()
	builtin[typ] = func(v any) lua.LGFunction {
		f := v.(func() (String, error))
		return func(state *lua.LState) int {
			v, err := f()
			if err != nil {
				state.RaiseError(err.Error())
				return 0
			}

			state.Push(lua.LString(v))
			return 1
		}
	}
}

func init() {
	typ := reflect.TypeOf((*func(Number) error)(nil)).Elem()
	builtin[typ] = func(v any) lua.LGFunction {
		f := v.(func(Number) error)
		return func(state *lua.LState) int {
			if err := f(Number(state.CheckNumber(1))); err != nil {
				state.RaiseError(err.Error())
			}
			return 0
		}
	}
}

func init() {
	typ := reflect.TypeOf((*func() (Number, error))(nil)).Elem()
	builtin[typ] = func(v any) lua.LGFunction {
		f := v.(func() (Number, error))
		return func(state *lua.LState) int {
			v, err := f()
			if err != nil {
				state.RaiseError(err.Error())
				return 0
			}

			state.Push(lua.LNumber(v))
			return 1
		}
	}
}

func init() {
	typ := reflect.TypeOf((*func(Bool) error)(nil)).Elem()
	builtin[typ] = func(v any) lua.LGFunction {
		f := v.(func(Bool) error)
		return func(state *lua.LState) int {
			if err := f(Bool(state.CheckBool(1))); err != nil {
				state.RaiseError(err.Error())
			}
			return 0
		}
	}
}

func init() {
	typ := reflect.TypeOf((*func() (Bool, error))(nil)).Elem()
	builtin[typ] = func(v any) lua.LGFunction {
		f := v.(func() (Bool, error))
		return func(state *lua.LState) int {
			v, err := f()
			if err != nil {
				state.RaiseError(err.Error())
				return 0
			}

			state.Push(lua.LBool(v))
			return 1
		}
	}
}
