// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

var (
	errInvalidScript = errors.New("script is not in a valid state")
)

// Script represents a LUA script
type Script struct {
	lock sync.Mutex
	name string         // The name of the script
	argn int            // The number of arguments
	exec *lua.LState    // The runtime for the script
	main *lua.LFunction // The main function
}

// FromReader reads a script fron an io.Reader
func FromReader(name string, r io.Reader) (*Script, error) {
	script := &Script{
		name: name,
	}
	return script, script.Update(r)
}

// FromString reads a script fron a string
func FromString(name, code string) (*Script, error) {
	return FromReader(name, bytes.NewBufferString(code))
}

// Run runs the main function of the script with arguments.
func (s *Script) Run(ctx context.Context, args ...interface{}) (Value, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.main == nil {
		return nil, errInvalidScript
	}

	// Push the arguments into the state
	exec := s.exec
	exec.SetContext(ctx)
	exec.Push(s.main)
	for _, arg := range args {
		exec.Push(luaValueOf(arg))

	}

	// Call the main function
	if err := exec.PCall(len(args), 1, nil); err != nil {
		return nil, err
	}

	// Pop the returned value
	result := exec.Get(-1)
	s.exec.Pop(1)
	return resultOf(result), nil
}

// Update updates the content of the script.
func (s *Script) Update(r io.Reader) error {
	runtime := newVM()
	fn, err := s.compile(r)
	if err != nil {
		return err
	}

	// Push the function to the runtime
	codeFn := runtime.NewFunctionFromProto(fn)
	runtime.Push(codeFn)

	// Initialize by calling the script
	if err := runtime.PCall(0, lua.MultRet, nil); err != nil {
		return err
	}

	// Get the main function
	mainFn, err := findFunction(runtime, "main")
	if err != nil {
		return err
	}

	// Make sure the most recent code is present in the state
	s.lock.Lock()
	defer s.lock.Unlock()
	s.argn = int(mainFn.Proto.NumParameters)
	s.exec = runtime
	s.main = mainFn
	return nil
}

// Compile compiles a script into a function that can be shared.
func (s *Script) compile(r io.Reader) (*lua.FunctionProto, error) {
	reader := bufio.NewReader(r)
	chunk, err := parse.Parse(reader, s.name)
	if err != nil {
		return nil, err
	}

	// Compile into a function
	proto, err := lua.Compile(chunk, s.name)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// Close closes the script and cleanly disposes of its resources.
func (s *Script) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.exec.Close()
	return nil
}

// newVM creates a new LUA state
func newVM() *lua.LState {
	state := lua.NewState(lua.Options{
		RegistrySize:        1024 * 20, // this is the initial size of the registry
		RegistryMaxSize:     1024 * 80, // this is the maximum size that the registry can grow to. If set to `0` (the default) then the registry will not auto grow
		RegistryGrowStep:    32,        // this is how much to step up the registry by each time it runs out of space. The default is `32`.
		CallStackSize:       120,       // this is the maximum callstack size of this LState
		MinimizeStackMemory: true,      // Defaults to `false` if not specified. If set, the callstack will auto grow and shrink as needed up to a max of `CallStackSize`. If not set, the callstack will be fixed at `CallStackSize`.
	})
	//state.PreloadModule("relay", vm.mod.loadModule)
	return state
}

// findFunction extracts a global function
func findFunction(runtime *lua.LState, name string) (*lua.LFunction, error) {
	fn := runtime.GetGlobal(name)
	if fn == nil {
		return nil, fmt.Errorf("%s() function not found", name)
	}

	if fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("%s() is %s, not a function", name, fn.Type())
	}

	return fn.(*lua.LFunction), nil
}
