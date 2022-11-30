// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestValue(typ Type) Value {
	switch typ {
	case TypeNumber:
		return Number(1)
	case TypeString:
		return String("x")
	case TypeBool:
		return Bool(true)
	default:
		return Nil{}
	}
}

func TestTableConvert(t *testing.T) {
	expect := newTestMap()
	input := ValueOf(expect)
	assert.EqualValues(t, expect, input.Native())
}

func TestComplexTable(t *testing.T) {
	expect := newComplexMap()
	input := ValueOf(expect)
	output := input.Native().([]any)
	assert.Len(t, output, len(expect))
}

// Test map[string][]float64
func TestTableCodec(t *testing.T) {
	expect := newTestMap()
	encoded, err := json.Marshal(expect)
	assert.NoError(t, err)

	// Decode a table
	var decoded Table
	assert.NoError(t, json.Unmarshal(encoded, &decoded))

	// re encode it back to json and compare strings
	reincoded, err := json.Marshal(decoded)
	assert.NoError(t, err)
	assert.Equal(t, encoded, reincoded)
}

func TestArrayMap(t *testing.T) {
	mp := []map[string]any{
		{"hello": []float64{1, 2, 3}},
		{"next": true},
		{"map": map[string]any{
			"a": []int64{1, 2},
			"b": 1,
			"c": false,
			"d": []map[string]any{
				{"e": "f", "g": "h"},
			},
		}},
		{"hello": "world"},
	}
	l := lua.NewState()
	val := lvalueOf(l, mp)
	tbl, ok := val.(*lua.LTable)
	assert.True(t, ok)

	mp1, ok := tbl.RawGetInt(1).(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, 3, mp1.RawGetString("hello").(*lua.LTable).Len())

	mp2, ok := tbl.RawGetInt(2).(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, lua.LTrue, mp2.RawGetString("next"))

	mp3, ok := tbl.RawGetInt(3).(*lua.LTable)
	assert.True(t, ok)
	mp3d, ok := mp3.RawGetString("map").(*lua.LTable).RawGetString("d").(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, 1, mp3d.Len())

	mp4, ok := tbl.RawGetInt(4).(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, lua.LString("world"), mp4.RawGetString("hello"))
}

func TestArraySlice(t *testing.T) {
	s := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	l := lua.NewState()
	val := lvalueOf(l, s)
	tbl, ok := val.(*lua.LTable)
	assert.True(t, ok)

	s1, ok := tbl.RawGetInt(1).(*lua.LTable)
	assert.Equal(t, 3, s1.Len())

	s1.ForEach(func(key, val lua.LValue) {
		assert.Equal(t, key, val)
	})
}

func TestNumbers(t *testing.T) {
	n := []int{1, 2, 3}
	l := lua.NewState()
	val := lvalueOf(l, n)
	tbl, ok := val.(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, 3, tbl.Len())

	tbl.ForEach(func(key, val lua.LValue) {
		assert.Equal(t, key, val)
	})
}

func TestBools(t *testing.T) {
	n := []bool{true, false, true}
	l := lua.NewState()
	val := lvalueOf(l, n)
	tbl, ok := val.(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, 3, tbl.Len())

	tbl.ForEach(func(key, val lua.LValue) {
		switch key {
		case lua.LNumber(1), lua.LNumber(3):
			assert.Equal(t, lua.LTrue, val)
		default:
			assert.Equal(t, lua.LFalse, val)
		}
	})
}

func TestStrings(t *testing.T) {
	n := []string{"aj", "roman", "abdo"}
	l := lua.NewState()
	val := lvalueOf(l, n)
	tbl, ok := val.(*lua.LTable)
	assert.True(t, ok)
	assert.Equal(t, 3, tbl.Len())

	tbl.ForEach(func(key, val lua.LValue) {
		switch key {
		case lua.LNumber(1):
			assert.Equal(t, lua.LString("aj"), val)
		case lua.LNumber(2):
			assert.Equal(t, lua.LString("roman"), val)
		case lua.LNumber(3):
			assert.Equal(t, lua.LString("abdo"), val)
		default:
			assert.Equal(t, lua.LFalse, val)
		}
	})
}

func TestMap(t *testing.T) {
	mp := map[string]any{
		"string":  "aj",
		"numbers": []float64{1, 2, 3},
		"bool":    true,
	}
	l := lua.NewState()
	val := lvalueOf(l, mp)
	tbl, ok := val.(*lua.LTable)
	assert.True(t, ok)
	tbl.ForEach(func(key, val lua.LValue) {
		switch key {
		case lua.LString("string"):
			assert.Equal(t, lua.LString("aj"), val)
		case lua.LString("numbers"):
			tbl1, ok := val.(*lua.LTable)
			assert.True(t, ok)
			assert.Equal(t, 3, tbl1.Len())
		case lua.LString("bool"):
			assert.Equal(t, lua.LTrue, val)
		}
	})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func newTestMap() map[string]any {
	return map[string]any{
		"user":   "Roman",
		"age":    37.0,
		"dev":    true,
		"bitmap": []bool{true, false, true},
		"floats": []float64{1, 2, 3, 4, 5},
	}
}

func newComplexMap() []map[string]any {
	return []map[string]any{
		{
			"args": []float64{2, 4},
		}, {
			"args": []float64{2, 4},
		},
	}
}
