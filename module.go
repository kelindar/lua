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

var builtin = make(map[reflect.Type]func(interface{}) lua.LGFunction, 8)

// Module represents a loadable module.
type Module interface {
	inject(state *lua.LState) error
}

// --------------------------------------------------------------------

// ScriptModule represents a loadable module written in LUA itself.
type ScriptModule struct {
	Script  *Script // The script that contains the module
	Name    string  // The name of the module
	Version string  // The module version string
}

// Inject loads the module into the state
func (m *ScriptModule) inject(runtime *lua.LState) error {

	// Inject the prerequisite modules of the module
	if err := m.Script.loadModules(runtime); err != nil {
		return err
	}

	// Push the function to the runtime
	codeFn := runtime.NewFunctionFromProto(m.Script.code)
	preload := runtime.GetField(runtime.GetField(runtime.Get(lua.EnvironIndex), "package"), "preload")
	if _, ok := preload.(*lua.LTable); !ok {
		return errors.New("package.preload must be a table")

	}
	runtime.SetField(preload, m.Name, codeFn)
	return nil
}

// --------------------------------------------------------------------

// NativeModule represents a loadable native module.
type NativeModule struct {
	lock    sync.Mutex
	funcs   map[string]fngen
	Name    string // The name of the module
	Version string // The module version string
}

type fngen struct {
	name string
	code interface{}
}

// Generate generates a function
func (g *fngen) generate() lua.LGFunction {
	rv := reflect.ValueOf(g.code)
	rt := rv.Type()
	if maker, ok := builtin[rt]; ok {
		return maker(g.code)
	}

	name := g.name
	argTypes := make([]Type, 0, rt.NumIn())
	for i := 0; i < rt.NumIn(); i++ {
		argTypes = append(argTypes, typeMap[rt.In(i)])
	}

	args := make([]reflect.Value, 0, rt.NumIn())
	return func(state *lua.LState) int {
		if state.GetTop() != len(argTypes) {
			state.RaiseError("%s expects %d arguments, but got %d", name, len(argTypes), state.GetTop())
			return 0
		}

		// Convert the arguments
		args = args[:0]
		for i, arg := range argTypes {
			switch arg {
			case TypeString:
				args = append(args, reflect.ValueOf(String(state.CheckString(i+1))))
			case TypeNumber:
				args = append(args, reflect.ValueOf(Number(state.CheckNumber(i+1))))
			case TypeBool:
				args = append(args, reflect.ValueOf(Bool(state.CheckBool(i+1))))
			case TypeStrings:
				args = append(args, reflect.ValueOf(resultOfStrings(state.CheckTable(i+1))))
			case TypeNumbers:
				args = append(args, reflect.ValueOf(resultOfNumbers(state.CheckTable(i+1))))
			case TypeBools:
				args = append(args, reflect.ValueOf(resultOfBools(state.CheckTable(i+1))))
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
	}
}

// Register registers a function into the module.
func (m *NativeModule) Register(name string, function interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Lazily create the function map
	if m.funcs == nil {
		m.funcs = make(map[string]fngen, 2)
	}

	// Validate the function
	if err := validate(name, function); err != nil {
		return err
	}

	m.funcs[name] = fngen{name: name, code: function}
	return nil
}

// Unregister unregisters a function from the module.
func (m *NativeModule) Unregister(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.funcs, name)
}

// Inject loads the module into the state
func (m *NativeModule) inject(state *lua.LState) error {
	table := make(map[string]lua.LGFunction, len(m.funcs))
	for name, g := range m.funcs {
		table[name] = g.generate()
	}

	state.PreloadModule(m.Name, func(state *lua.LState) int {
		mod := state.SetFuncs(state.NewTable(), table)
		state.SetField(mod, "version", lua.LString(m.Version))
		state.Push(mod)
		return 1
	})
	return nil
}

// validate validates the function type
func validate(name string, function interface{}) error {
	rv := reflect.ValueOf(function)
	rt := rv.Type()
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
