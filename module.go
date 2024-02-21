// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	errFuncInput   = errors.New("lua: function input arguments must be of type lua.Value")
	errFuncOutput  = errors.New("lua: function return values must be either (error) or (lua.Value, error)")
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

var builtin = make(map[reflect.Type]func(any) lua.LGFunction, 8)

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
	code any
}

// Generate generates a function
func (g *fngen) generate() lua.LGFunction {
	rv := reflect.ValueOf(g.code)
	rt := rv.Type()
	if maker, ok := builtin[rt]; ok {
		return maker(g.code)
	}

	name := g.name
	args := make([]reflect.Value, 0, rt.NumIn())
	return func(state *lua.LState) int {
		if state.GetTop() != rt.NumIn() {
			state.RaiseError("%s expects %d arguments, but got %d", name, rt.NumIn(), state.GetTop())
			return 0
		}

		// Convert the arguments
		args = args[:0]
		for i := 0; i < rt.NumIn(); i++ {
			args = append(args, reflect.ValueOf(resultOf(state.Get(i+1))))
		}

		// Call the function
		out := rv.Call(args)
		switch len(out) {
		case 1:
			if err := out[0]; !err.IsNil() {
				state.RaiseError(err.Interface().(error).Error())
			}
			return 0
		default:
			if err := out[1]; !err.IsNil() {
				state.RaiseError(err.Interface().(error).Error())
				return 0
			}
			state.Push(lvalueOf(state, out[0].Interface()))
			return 1
		}
	}
}

// Register registers a function into the module.
func (m *NativeModule) Register(name string, function any) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Lazily create the function map
	if m.funcs == nil {
		m.funcs = make(map[string]fngen, 2)
	}

	// Validate the function
	if err := validate(function); err != nil {
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
func validate(function any) error {
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
	case rt.NumOut() == 1 && isError(rt, 0):
	case rt.NumOut() == 2 && isValid(rt, 0) && isError(rt, 1):
	default:
		return errFuncOutput
	}
	return nil
}

func isError(rt reflect.Type, at int) bool {
	return rt.Out(at).Implements(typeError)
}

func isValid(rt reflect.Type, at int) bool {
	switch rt.Out(at) {
	case typeString:
	case typeNumber:
	case typeBool:
	case typeNumbers:
	case typeStrings:
	case typeBools:
	case typeTable:
	case typeArray:
	case typeValue:
	default:
		return false
	}
	return true
}
