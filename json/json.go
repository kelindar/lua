// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

// This is a fork of https://github.com/layeh/gopher-json, licensed under The Unlicense

package json

import (
	"encoding/json"
	"errors"

	lua "github.com/yuin/gopher-lua"
)

var (
	errNested      = errors.New("cannot encode recursively nested tables to JSON")
	errSparseArray = errors.New("cannot encode sparse array")
	errInvalidKeys = errors.New("cannot encode mixed or invalid key types")
)

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"decode": apiDecode,
	"encode": apiEncode,
	"array":  apiArray,
}

func apiDecode(state *lua.LState) int {
	str := state.CheckString(1)

	value, err := Decode(state, []byte(str))
	if err != nil {
		state.Push(lua.LNil)
		state.Push(lua.LString(err.Error()))
		return 2
	}
	state.Push(value)
	return 1
}

func apiEncode(state *lua.LState) int {
	value := state.CheckAny(1)

	data, err := Encode(value)
	if err != nil {
		state.Push(lua.LNil)
		state.RaiseError(err.Error())
		return 1
	}
	state.Push(lua.LString(string(data)))
	return 1
}

// --------------------------------------------------------------------

// EmptyArray is a marker for an empty array.
var EmptyArray = &lua.LUserData{Value: []any(nil)}

// apiArray creates an array from the arguments.
func apiArray(state *lua.LState) int {
	switch state.GetTop() {
	case 0:
		state.Push(EmptyArray)
		return 1
	case 1:
		// If it's not a table, or empty return a marker
		table, ok := state.CheckAny(1).(*lua.LTable)
		switch {
		case ok && table.Len() > 0: // array
			state.Push(table)
			return 1
		case ok: // check if it's an empty map
			k, v := table.Next(lua.LNil)
			if k == lua.LNil && v == lua.LNil {
				state.Push(EmptyArray)
				return 1
			}
		}

		// Otherwise, return an array
		fallthrough
	default:
		table := state.CreateTable(state.GetTop(), 0)
		for i := 1; i <= state.GetTop(); i++ {
			table.RawSetInt(i, state.Get(i))
		}

		// Return the table
		state.Push(table)
		return 1
	}
}

// --------------------------------------------------------------------

type invalidTypeError lua.LValueType

func (i invalidTypeError) Error() string {
	return `cannot encode ` + lua.LValueType(i).String() + ` to JSON`
}

// Encode returns the JSON encoding of value.
func Encode(value lua.LValue) ([]byte, error) {
	return json.Marshal(jsonValue{
		LValue:  value,
		visited: make(map[*lua.LTable]bool),
	})
}

type jsonValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (j jsonValue) MarshalJSON() (data []byte, err error) {
	switch converted := j.LValue.(type) {
	case lua.LBool:
		data, err = json.Marshal(bool(converted))
	case lua.LNumber:
		data, err = json.Marshal(float64(converted))
	case *lua.LNilType:
		data = []byte(`null`)
	case *lua.LUserData:
		switch {
		case converted == EmptyArray:
			data = []byte(`[]`)
		default:
			data, err = json.Marshal(converted.Value)
		}
	case lua.LString:
		data, err = json.Marshal(string(converted))
	case *lua.LTable:
		if j.visited[converted] {
			return nil, errNested
		}
		j.visited[converted] = true

		key, value := converted.Next(lua.LNil)

		switch key.Type() {
		case lua.LTNil: // empty table
			data = []byte(`[]`)
		case lua.LTNumber:
			arr := make([]jsonValue, 0, converted.Len())
			expectedKey := lua.LNumber(1)
			for key != lua.LNil {
				if key.Type() != lua.LTNumber {
					err = errInvalidKeys
					return
				}
				if expectedKey != key {
					err = errSparseArray
					return
				}
				arr = append(arr, jsonValue{value, j.visited})
				expectedKey++
				key, value = converted.Next(key)
			}
			data, err = json.Marshal(arr)
		case lua.LTString:
			obj := make(map[string]jsonValue)
			for key != lua.LNil {
				if key.Type() != lua.LTString {
					err = errInvalidKeys
					return
				}
				obj[key.String()] = jsonValue{value, j.visited}
				key, value = converted.Next(key)
			}
			data, err = json.Marshal(obj)
		default:
			err = errInvalidKeys
		}
	default:
		err = invalidTypeError(j.LValue.Type())
	}
	return
}

// Decode converts the JSON encoded data to Lua values.
func Decode(L *lua.LState, data []byte) (lua.LValue, error) {
	var value any
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}
	return DecodeValue(L, value), nil
}

// DecodeValue converts the value to a Lua value.
//
// This function only converts values that the encoding/json package decodes to.
// All other values will return lua.LNil.
func DecodeValue(L *lua.LState, value any) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case json.Number:
		return lua.LString(converted)
	case []any:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(DecodeValue(L, item))
		}
		return arr
	case map[string]any:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), DecodeValue(L, item))
		}
		return tbl
	case nil:
		return lua.LNil
	}

	return lua.LNil
}
