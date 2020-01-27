// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/yuin/gopher-lua"
)

var (
	errFuncInput  = errors.New("lua: function input arguments must be of type lua.Value")
	errFuncOutput = errors.New("lua: function return values must be either (error) or (lua.Value, error)")
)

// Module represents a loadable module
type Module struct {
	lock    sync.Mutex
	funcs   map[string]lua.LGFunction
	Name    string // The name of the module
	Version string // The module version string
}

// Generate generates a function
func generate(name string, function interface{}) (lua.LGFunction, error) {
	rv := reflect.ValueOf(function)
	rt := rv.Type()
	if err := isValidFunction(rt); err != nil {
		return nil, err
	}

	argCount := rt.NumIn()
	argTypes := make([]Type, 0, rt.NumIn())
	for i := 0; i < rt.NumIn(); i++ {
		argTypes = append(argTypes, typeMap[rt.In(i)])
	}

	return func(state *lua.LState) int {
		if state.GetTop() != argCount {
			state.RaiseError("%s expects %d arguments, but got %d", name, argCount, state.GetTop())
			return 0
		}

		// Convert the arguments
		args := make([]reflect.Value, 0, argCount)
		for i, arg := range argTypes {
			switch arg {
			case TypeString:
				args = append(args, reflect.ValueOf(String(state.CheckString(i+1))))
			case TypeNumber:
				args = append(args, reflect.ValueOf(Number(state.CheckNumber(i+1))))
			case TypeBool:
				args = append(args, reflect.ValueOf(Bool(state.CheckBool(i+1))))
			}
		}

		// Call the function
		out := rv.Call(args)
		switch len(out) {
		case 1:
			if err := out[0]; !err.IsNil() {
				state.RaiseError(err.String())
			}
			return 0
		default:
			if err := out[1]; !err.IsNil() {
				state.RaiseError(err.String())
				return 0
			}

			state.Push(luaValueOf(out[0]))
			return 1
		}
	}, nil
}

// Register registers a function into the module.
func (m *Module) Register(name string, function interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Lazily create the function map
	if m.funcs == nil {
		m.funcs = make(map[string]lua.LGFunction, 2)
	}

	// Generate the function
	f, err := generate(name, function)
	if err != nil {
		return err
	}

	m.funcs[name] = f
	return nil
}

// Unregister unregisters a function from the module.
func (m *Module) Unregister(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.funcs, name)
}

// Inject loads the module into the state
func (m *Module) inject(state *lua.LState) {
	state.PreloadModule(m.Name, func(state *lua.LState) int {
		mod := state.SetFuncs(state.NewTable(), m.funcs)
		state.SetField(mod, "version", lua.LString(m.Version))
		state.Push(mod)
		return 1
	})
}

// isValidFunction validates the function type
func isValidFunction(rt reflect.Type) error {
	if rt.Kind() != reflect.Func {
		return fmt.Errorf("lua: input is a %s, not a function", rt.Kind().String())
	}

	// Validate the input
	for i := 0; i < rt.NumIn(); i++ {
		if _, ok := typeMap[rt.In(i)]; !ok {
			return errFuncInput
		}
	}

	// Validate the output
	switch {
	case rt.NumOut() == 1 && rt.Out(0).Implements(typeError):
	case rt.NumOut() == 2 && rt.Out(0) == typeString && rt.Out(1).Implements(typeError):
	case rt.NumOut() == 2 && rt.Out(0) == typeNumber && rt.Out(1).Implements(typeError):
	case rt.NumOut() == 2 && rt.Out(0) == typeBool && rt.Out(1).Implements(typeError):
	default:
		return errFuncOutput
	}
	return nil
}
