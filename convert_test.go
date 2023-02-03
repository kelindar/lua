// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestValueOf(t *testing.T) {
	tests := []struct {
		input  any
		output Value
	}{
		{input: complex64(1), output: Nil{}},
		{input: Number(1), output: Number(1)},
		{input: Bool(true), output: Bool(true)},
		{input: String("hi"), output: String("hi")},
		{input: int(1), output: Number(1)},
		{input: int8(1), output: Number(1)},
		{input: int16(1), output: Number(1)},
		{input: int32(1), output: Number(1)},
		{input: int64(1), output: Number(1)},
		{input: uint(1), output: Number(1)},
		{input: uint8(1), output: Number(1)},
		{input: uint16(1), output: Number(1)},
		{input: uint32(1), output: Number(1)},
		{input: uint64(1), output: Number(1)},
		{input: float32(1), output: Number(1)},
		{input: float64(1), output: Number(1)},
		{input: true, output: Bool(true)},
		{input: "hi", output: String("hi")},
		{input: []string{"a", "b"}, output: Strings{"a", "b"}},
		{input: []int{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []int8{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []int16{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []int32{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []int64{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []uint{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []uint8{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []uint16{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []uint32{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []uint64{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []float32{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []float64{1, 2, 3}, output: Numbers{1, 2, 3}},
		{input: []bool{false, true}, output: Bools{false, true}},
		{input: nil, output: Nil{}},
		{input: Nil{}, output: Nil{}},
		{input: map[string]any{
			"A": "foo",
			"B": "bar",
		}, output: Table{
			"A": String("foo"),
			"B": String("bar"),
		}},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.output, ValueOf(tc.input))
		assert.NotEmpty(t, tc.output.String())
	}
}

func TestResultOfArray(t *testing.T) {
	mp := []map[string]any{
		{"hello": []float64{1, 2, 3}},
		{"next": true},
		{"hello": "world"},
	}
	l := lua.NewState()
	val := lvalueOf(l, mp)
	res := resultOf(val)
	array, ok := res.(Array)
	assert.True(t, ok)
	assert.Len(t, array, 3)

	resMp := array.Native().([]any)
	for i, val := range resMp {
		assert.Equal(t, mp[i], val.(map[string]any))
	}
}

func TestResultComplexMap(t *testing.T) {
	mp := []map[string]any{
		{"map": []map[string]any{
			{"e": "f", "g": "h"}},
		},
	}
	l := lua.NewState()
	val := lvalueOf(l, mp)
	res := resultOf(val)
	array, ok := res.(Array)
	assert.True(t, ok)
	resMp := array.Native().([]any)
	assert.Len(t, resMp, len(mp))
}

func TestResultOfMap(t *testing.T) {
	mp := map[string]any{
		"string":  "aj",
		"numbers": []float64{1, 2, 3},
		"bool":    true,
	}
	l := lua.NewState()
	val := lvalueOf(l, mp)
	res := resultOf(val)
	table, ok := res.(Table)
	assert.True(t, ok)
	assert.Len(t, table, 3)
	resMp, ok := table.Native().(map[string]any)
	assert.True(t, ok)

	for _, key := range []string{"string", "numbers", "bool"} {
		assert.Equal(t, mp[key], resMp[key])
	}
}
