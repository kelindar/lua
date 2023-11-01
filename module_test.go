// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"context"
	"errors"
	"fmt"
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
	must(m.Register("enrich", enrich))
	must(m.Register("error", errorfunc))
	must(m.Register("error1", errorfunc1))
	must(m.Register("toNumbers", toNumbers))
	return m
}

func toNumbers(v Value) (Numbers, error) {
	switch t := v.(type) {
	case Numbers:
		return t, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}

func sum(a, b Number) (Number, error) {
	return a + b, nil
}

func echo(v String) (String, error) {
	return v, nil
}

func errorfunc(v String) (String, error) {
	return "", fmt.Errorf("error with input (%v)", v)
}

func errorfunc1(_ Table) (String, error) {
	return "", errors.New("throwing error")
}

func hash(s Value) (Number, error) {
	h := fnv.New32a()
	h.Write([]byte(s.(String)))

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

func enrich(name String, request Table) (Table, error) {
	request["name"] = name
	request["age"] = Number(30)
	return request, nil
}

func toNumbers(v Value) (Numbers, error) {
	switch t := v.(type) {
	case Numbers:
		return t, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
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

func Test_Value(t *testing.T) {
	s, err := newScript("fixtures/value.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), []float64{1.1, 2.1})
	assert.NoError(t, err)
	assert.Equal(t, TypeNumbers, out.Type())
	assert.Equal(t, []float64{1.1, 2.1}, out.Native())
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

	out, err := s.Run(context.Background(), map[string]any{
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

func TestEnrich(t *testing.T) {
	s, err := newScript("fixtures/enrich.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), map[string]any{
		"A": "apples",
		"B": "oranges",
	})
	assert.NoError(t, err)
	assert.Equal(t, TypeTable, out.Type())
	assert.EqualValues(t, map[string]any{
		"A":    "apples",
		"B":    "oranges",
		"age":  30.0,
		"name": "roman",
	}, out.(Table).Native())
}

func TestErrorMessage(t *testing.T) {
	s, err := newScript("fixtures/error.lua")
	assert.NoError(t, err)

	_, err = s.Run(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error with input (roman)")
}

func TestErrorMessage1(t *testing.T) {
	s, err := newScript("fixtures/error1.lua")
	assert.NoError(t, err)

	_, err = s.Run(context.Background(), map[string]string{
		"test1": "default",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "throwing error")
}

func TestEnrichComplexTable(t *testing.T) {
	s, err := newScript("fixtures/enrich.lua")
	assert.NoError(t, err)

	v, err := s.Run(context.Background(), map[string][]float64{
		"A": {1, 2, 3},
		"B": {1, 2, 3},
	})

	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, Table{
		"A":    Numbers{1, 2, 3},
		"B":    Numbers{1, 2, 3},
		"age":  Number(30),
		"name": String("roman"),
	}, v)
}

func TestUserdata(t *testing.T) {
	s, err := FromString("sandbox", `
	function main(request)
		return type(request)
    end
	`)
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), []map[string]any{
		{"a": 1.0, "b": 2.0},
		{"a": 10.0, "b": 20.0},
	})

	assert.NoError(t, err)
	assert.Equal(t, String("table"), out)
}

func Test_Any(t *testing.T) {
	s, err := newScript("fixtures/any.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), []float64{1.1, 2.1})
	assert.NoError(t, err)
	assert.Equal(t, TypeNumbers, out.Type())
	assert.Equal(t, []float64{1.1, 2.1}, out.Native())
}
