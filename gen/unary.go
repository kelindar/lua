package lua

import (
	"github.com/yuin/gopher-lua"
	"reflect"
)

//go:generate genny -in=$GOFILE -out=../z_unary.go gen "TIn=String,Number,Bool"

func init() {
	typ := reflect.TypeOf((*func(TIn) error)(nil)).Elem()
	builtin[typ] = func(v interface{}) lua.LGFunction {
		f := v.(func(TIn) error)
		return func(state *lua.LState) int {
			if err := f(TIn(state.CheckTIn(1))); err != nil {
				state.RaiseError(err.Error())
			}
			return 0
		}
	}
}

func init() {
	typ := reflect.TypeOf((*func() (TIn, error))(nil)).Elem()
	builtin[typ] = func(v interface{}) lua.LGFunction {
		f := v.(func() (TIn, error))
		return func(state *lua.LState) int {
			v, err := f()
			if err != nil {
				state.RaiseError(err.Error())
				return 0
			}

			state.Push(lua.LTIn(v))
			return 1
		}
	}
}
