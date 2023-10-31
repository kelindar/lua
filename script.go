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
	"math"
	"runtime"
	"sync"

	"github.com/kelindar/lua/json"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

var (
	errInvalidScript = errors.New("lua: script is not in a valid state")
)

// Script represents a LUA script
type Script struct {
	lock sync.RWMutex
	name string             // The name of the script
	conc int                // The concurrency setting for the VM pool
	pool pool               // The pool of runtimes for concurrent use
	mods []Module           // The injected modules
	code *lua.FunctionProto // The precompiled code
}

// New creates a new script from an io.Reader
func New(name string, source io.Reader, concurrency int, modules ...Module) (*Script, error) {
	if concurrency <= 0 {
		concurrency = defaultConcurrency
	}

	script := &Script{
		name: name,
		mods: modules,
		conc: concurrency,
	}
	return script, script.Update(source)
}

// FromReader reads a script fron an io.Reader
func FromReader(name string, r io.Reader, modules ...Module) (*Script, error) {
	return New(name, r, 0, modules...)
}

// FromString reads a script fron a string
func FromString(name, code string, modules ...Module) (*Script, error) {
	return New(name, bytes.NewBufferString(code), 0, modules...)
}

// Name returns the name of the script
func (s *Script) Name() string {
	return s.name
}

// Concurrency returns the concurrency setting of the script
func (s *Script) Concurrency() int {
	return s.conc
}

// Run runs the main function of the script with arguments.
func (s *Script) Run(ctx context.Context, args ...any) (Value, error) {

	// Protect swapping of the pools when the script is updated.
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Acquire and release the pool of VMs, given our read lock we can still
	// enter here concurrently so the pool must also be thread-safe.
	vm := s.pool.Acquire()
	defer s.pool.Release(vm)

	// Run the script
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
	s.pool, err = newPool(s, s.conc)
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
	tbls *tablePool     // Pool of tables
}

// newVM creates a new VM for a script
func newVM(s *Script) (*vm, error) {
	l := newState()
	v := &vm{
		tbls: newTablePool(l),
		exec: l,
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
func (v *vm) Run(ctx context.Context, args []any) (Value, error) {
	if v.main == nil {
		return nil, errInvalidScript
	}

	// Push the arguments into the state
	exec := v.exec
	exec.SetContext(ctx)
	exec.Push(v.main)
	for _, arg := range args {
		switch x := arg.(type) {
		case Table:
			defer pushArg(x, exec, v.tbls)()
		case Bools:
			defer pushArg(x, exec, v.tbls)()
		case Strings:
			defer pushArg(x, exec, v.tbls)()
		case Numbers:
			defer pushArg(x, exec, v.tbls)()
		case Array:
			defer pushArg(x, exec, v.tbls)()
		default:
			exec.Push(lvalueOf(exec, arg))
		}
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
		RegistrySize:        1024 * 4,   // this is the initial size of the registry
		RegistryMaxSize:     1024 * 128, // this is the maximum size that the registry can grow to. If set to `0` (the default) then the registry will not auto grow
		RegistryGrowStep:    32,         // this is how much to step up the registry by each time it runs out of space. The default is `32`.
		CallStackSize:       64,         // this is the maximum callstack size of this LState
		MinimizeStackMemory: true,       // Defaults to `false` if not specified. If set, the callstack will auto grow and shrink as needed up to a max of `CallStackSize`. If not set, the callstack will be fixed at `CallStackSize`.
	})
}

// --------------------------------------------------------------------

type tablePool struct {
	sync.Pool
}

// newTablePool creates a new table pool
func newTablePool(state *lua.LState) *tablePool {
	return &tablePool{
		Pool: sync.Pool{
			New: func() any { return state.NewTable() },
		},
	}
}

// Get gets a table from the pool
func (b *tablePool) Get() *lua.LTable {
	return b.Pool.Get().(*lua.LTable)
}

// Put puts a table back into the pool
func (b *tablePool) Put(x *lua.LTable) {
	b.Pool.Put(x)
}

// pushArg pushes an argument into the state
func pushArg[T tableType](value T, state *lua.LState, p *tablePool) func() {
	dst := p.Get()
	value.lcopy(dst, state)
	state.Push(dst)
	return func() {
		p.Put(dst)
	}
}

// --------------------------------------------------------------------

// defaultConcurrency sets the default concurrency for the VM pool
var defaultConcurrency = int(math.Min(
	float64(runtime.GOMAXPROCS(-1)), float64(runtime.NumCPU()),
))

// Pool holds a pool of runtimes.
type pool chan *vm

// newPool creates a new pool of runtimes.
func newPool(s *Script, concurrency int) (pool, error) {
	pool := make(pool, concurrency)
	for i := 0; i < concurrency; i++ {
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
	return <-p // Wait until we have a VM
}

// Release returns a state to the pool.
func (p pool) Release(vm *vm) {
	select {
	case p <- vm:
	default: // Discard
	}
}
