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
	"runtime"
	"sync"

	"github.com/kelindar/lua/json"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"layeh.com/gopher-luar"
)

var (
	errInvalidScript = errors.New("lua: script is not in a valid state")
)

// Script represents a LUA script
type Script struct {
	lock sync.RWMutex
	name string             // The name of the script
	pool pool               // The pool of runtimes for concurrent use
	mods []Module           // The injected modules
	code *lua.FunctionProto // The precompiled code
}

// FromReader reads a script fron an io.Reader
func FromReader(name string, r io.Reader, modules ...Module) (*Script, error) {
	script := &Script{
		name: name,
		mods: modules,
	}
	return script, script.Update(r)
}

// FromString reads a script fron a string
func FromString(name, code string, modules ...Module) (*Script, error) {
	return FromReader(name, bytes.NewBufferString(code), modules...)
}

// Name returns the name of the script
func (s *Script) Name() string {
	return s.name
}

// Run runs the main function of the script with arguments.
func (s *Script) Run(ctx context.Context, args ...interface{}) (Value, error) {

	// Protect swapping of the pools when the script is updated.
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Acquire and release the pool of VMs, given our read lock we can still
	// enter here concurrently so the pool must also be thread-safe.
	vm := s.pool.Acquire()
	if vm == nil {
		return nil, errInvalidScript
	}

	defer s.pool.Release(vm)
	return vm.Run(ctx, args)
}

// Update updates the content of the script.
func (s *Script) Update(r io.Reader) (err error) {
	code, err := s.compile(r)
	if err != nil {
		return err
	}

	// Protect from now on, as we need to update the script while loading
	s.lock.Lock()
	defer s.lock.Unlock()

	// Create a new pool of VMs
	s.code = code
	s.pool, err = newPool(s)
	return
}

// Compile compiles a script into a function that can be shared.
func (s *Script) compile(r io.Reader) (*lua.FunctionProto, error) {
	reader := bufio.NewReader(r)
	chunk, err := parse.Parse(reader, s.name)
	if err != nil {
		return nil, err
	}

	// Compile into a function
	return lua.Compile(chunk, s.name)
}

// LoadModules loads in the prerequisite modules
func (s *Script) loadModules(runtime *lua.LState) error {
	runtime.PreloadModule("json", json.Loader)
	for _, m := range s.mods {
		if err := m.inject(runtime); err != nil {
			return err
		}
	}

	return nil
}

// Close closes the script and cleanly disposes of its resources.
func (s *Script) Close() error {
	return nil
}

// findFunction extracts a global function
func findFunction(runtime *lua.LState, name string) (*lua.LFunction, error) {
	fn := runtime.GetGlobal(name)
	if fn == nil || fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("lua: %s() function not found", name)
	}

	return fn.(*lua.LFunction), nil
}

// --------------------------------------------------------------------

// VM represents a single VM which can only be ran serially.
type vm struct {
	argn int            // The number of arguments
	exec *lua.LState    // The pool of runtimes for concurrent use
	main *lua.LFunction // The main function
}

// newVM creates a new VM for a script
func newVM(s *Script) (*vm, error) {
	v := &vm{
		exec: newState(),
		argn: 0,
		main: nil,
	}

	// Push the function to the runtime
	codeFn := v.exec.NewFunctionFromProto(s.code)
	v.exec.Push(codeFn)

	// Inject the modules
	if err := s.loadModules(v.exec); err != nil {
		return nil, err
	}

	// Initialize by calling the script
	if err := v.exec.PCall(0, lua.MultRet, nil); err != nil {
		return nil, err
	}

	// If we have a main function, set it
	if mainFn, err := findFunction(v.exec, "main"); err == nil {
		v.argn = int(mainFn.Proto.NumParameters)
		v.main = mainFn
	}
	return v, nil
}

// Run runs the main function of the script with arguments.
func (v *vm) Run(ctx context.Context, args []interface{}) (Value, error) {
	if v.main == nil {
		return nil, errInvalidScript
	}

	// Push the arguments into the state
	exec := v.exec
	exec.SetContext(ctx)
	exec.Push(v.main)
	for _, arg := range args {
		exec.Push(luar.New(exec, arg))
	}

	// Call the main function
	if err := exec.PCall(len(args), 1, nil); err != nil {
		return nil, err
	}

	// Pop the returned value
	result := exec.Get(-1)
	exec.Pop(1)
	return resultOf(result), nil
}

// newState creates a new LUA state
func newState() *lua.LState {
	return lua.NewState(lua.Options{
		RegistrySize:        1024 * 20, // this is the initial size of the registry
		RegistryMaxSize:     1024 * 80, // this is the maximum size that the registry can grow to. If set to `0` (the default) then the registry will not auto grow
		RegistryGrowStep:    32,        // this is how much to step up the registry by each time it runs out of space. The default is `32`.
		CallStackSize:       120,       // this is the maximum callstack size of this LState
		MinimizeStackMemory: true,      // Defaults to `false` if not specified. If set, the callstack will auto grow and shrink as needed up to a max of `CallStackSize`. If not set, the callstack will be fixed at `CallStackSize`.
	})
}

// --------------------------------------------------------------------

// Pool holds a pool of runtimes.
type pool chan *vm

// newPool creates a new pool of runtimes.
func newPool(s *Script) (pool, error) {
	max := runtime.NumCPU()
	pool := make(pool, max)
	for i := 0; i < max; i++ {
		vm, err := newVM(s)
		if err != nil {
			return nil, err
		}

		pool <- vm
	}

	return pool, nil
}

// Acquire gets a state from the pool.
func (p pool) Acquire() (vm *vm) {
	select {
	case vm = <-p:
	default:
		vm = nil
	}
	return
}

// Release returns a state to the pool.
func (p pool) Release(vm *vm) {
	select {
	case p <- vm:
	default: // Discard
	}
}
