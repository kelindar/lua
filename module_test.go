// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"context"
	"hash/fnv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testModule() Module {
	m := &NativeModule{
		Name:    "test",
		Version: "1.0.0",
	}
	must(m.Register("hash", hash))
	must(m.Register("echo", echo))
	must(m.Register("sum", sum))
	must(m.Register("join", join))
	must(m.Register("sleep", sleep))
	must(m.Register("joinMap", joinMap))
	return m
}

func sum(a, b Number) (Number, error) {
	return a + b, nil
}

func echo(v String) (String, error) {
	return v, nil
}

func hash(s String) (Number, error) {
	h := fnv.New32a()
	h.Write([]byte(s))

	return Number(h.Sum32()), nil
}

func join(v Strings) (String, error) {
	return String(strings.Join([]string(v), ", ")), nil
}

func sleep(v Number) error {
	time.Sleep(time.Duration(v) * time.Millisecond)
	return nil
}

func joinMap(table Table) (String, error) {
	var sb strings.Builder
	for k, v := range table {
		sb.WriteString(k + ": " + v.String() + ", ")
	}
	return String(sb.String()), nil
}

func Test_Join(t *testing.T) {
	s, err := newScript("fixtures/join.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, TypeString, out.Type())
	assert.Equal(t, "apples, oranges, watermelons", string(out.(String)))
}

func Test_Hash(t *testing.T) {
	s, err := newScript("fixtures/hash.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), "abcdef")
	assert.NoError(t, err)
	assert.Equal(t, TypeNumber, out.Type())
	assert.Equal(t, int64(4282878506), int64(out.(Number)))
}

func Test_Sum(t *testing.T) {
	s, err := newScript("fixtures/sum.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, TypeNumber, out.Type())
	assert.Equal(t, int64(5), int64(out.(Number)))
}

func Test_JoinMap(t *testing.T) {
	s, err := newScript("fixtures/joinMap.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), map[string]interface{}{
		"A": "apples",
		"B": "oranges",
	})
	assert.NoError(t, err)
	assert.Equal(t, TypeString, out.Type())
	assert.Contains(t, string(out.(String)), "A: apples")
	assert.Contains(t, string(out.(String)), "B: oranges")
}

func Test_NotAFunc(t *testing.T) {
	m := &NativeModule{
		Name:    "test",
		Version: "1.0.0",
	}
	assert.Error(t, m.Register("xxx", 123))
	assert.NoError(t, m.Register("hash", hash))
	assert.Equal(t, 1, len(m.funcs))
	m.Unregister("hash")
	assert.Equal(t, 0, len(m.funcs))
}

func Test_ScriptModule(t *testing.T) {

	m, err := newScript("fixtures/module.lua")
	assert.NoError(t, err)

	s, err := newScript("fixtures/demo.lua", &ScriptModule{
		Script:  m,
		Name:    "demo_mod",
		Version: "1.0.0",
	})
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), 10, m)
	assert.NoError(t, err)
	assert.Equal(t, TypeNumber, out.Type())
	assert.Equal(t, Number(25), out.(Number))
	assert.Equal(t, "25", out.String())

	err = s.Close()
	assert.NoError(t, err)
}
