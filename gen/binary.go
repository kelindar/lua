package lua

import (
	"github.com/cheekybits/genny/generic"
	"github.com/yuin/gopher-lua"
	"reflect"
)

//go:generate genny -in=$GOFILE -out=../z_binary.go gen "TIn=String,Number,Bool TOut=String,Number,Bool"

// TIn is the generic type.
type TIn generic.Type

// TOut is the generic type.
type TOut generic.Type

func init() {
	typ := reflect.TypeOf((*func(TIn) (TOut, error))(nil)).Elem()
	builtin[typ] = func(v interface{}) lua.LGFunction {
		f := v.(func(TIn) (TOut, error))
		return func(state *lua.LState) int {
			v, err := f(TIn(state.CheckTIn(1)))
			if err != nil {
				state.RaiseError(err.Error())
				return 0
			}

			state.Push(luaValueOf(v))
			return 1
		}
	}
}
