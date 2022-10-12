// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"encoding/json"
	"reflect"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// ValueOf converts the type to our value
func ValueOf(v any) Value {
	switch v := v.(type) {
	case Number:
		return v
	case String:
		return v
	case Bool:
		return v
	case Numbers:
		return v
	case Strings:
		return v
	case Bools:
		return v
	case Table:
		return v
	case Tables:
		return v
	case int:
		return Number(v)
	case int8:
		return Number(v)
	case int16:
		return Number(v)
	case int32:
		return Number(v)
	case int64:
		return Number(v)
	case uint:
		return Number(v)
	case uint8:
		return Number(v)
	case uint16:
		return Number(v)
	case uint32:
		return Number(v)
	case uint64:
		return Number(v)
	case float32:
		return Number(v)
	case float64:
		return Number(v)
	case bool:
		return Bool(v)
	case string:
		return String(v)
	case []int:
		return numbersOf(v)
	case []int8:
		return numbersOf(v)
	case []int16:
		return numbersOf(v)
	case []int32:
		return numbersOf(v)
	case []int64:
		return numbersOf(v)
	case []uint:
		return numbersOf(v)
	case []uint8:
		return numbersOf(v)
	case []uint16:
		return numbersOf(v)
	case []uint32:
		return numbersOf(v)
	case []uint64:
		return numbersOf(v)
	case []float32:
		return numbersOf(v)
	case []float64:
		return Numbers(v)
	case []bool:
		return Bools(v)
	case []string:
		return Strings(v)
	case map[string]any:
		return mapAsTable(v)
	case []map[string]any:
		return mapsAsTables(v)
	case []any:
		return asArray(v)
	case nil:
		return Nil{}
	default:
		out, err := json.Marshal(v)
		if err != nil {
			return Nil{}
		}

		var resp any
		if err := json.Unmarshal(out, &resp); err != nil {
			return Nil{}
		}

		return ValueOf(resp)
	}
}

// ValueOf converts a value to a LUA-friendly one.
func resultOf(v lua.LValue) Value {
	switch v := v.(type) {
	case lua.LNumber:
		return Number(v)
	case lua.LString:
		return String(v)
	case lua.LBool:
		return Bool(v)
	case *lua.LTable:
		if top := v.RawGetInt(1); top != nil {
			switch top.Type() {
			case lua.LTNil:
				return asTable(v)
			case lua.LTNumber:
				return asNumbers(v)
			case lua.LTString:
				return asStrings(v)
			case lua.LTBool:
				return asBools(v)
			case lua.LTTable:
				return asTables(v)
			}
		}
		return Nil{}
	case *lua.LUserData:
		return ValueOf(v.Value)
	default:
		return Nil{}
	}
}

func asNumbers(t *lua.LTable) (out Numbers) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, float64(v.(lua.LNumber)))
	})
	return
}

func asStrings(t *lua.LTable) (out Strings) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, string(v.(lua.LString)))
	})
	return
}

func asBools(t *lua.LTable) (out Bools) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, bool(v.(lua.LBool)))
	})
	return
}

func asTable(t *lua.LTable) Table {
	out := make(Table, t.Len())
	t.ForEach(func(k, v lua.LValue) {
		out[k.String()] = resultOf(v)
	})
	return out
}

func asTables(t *lua.LTable) (out Tables) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, asTable(v.(*lua.LTable)))
	})
	return
}

func asArray(input []any) Array {
	arr := make(Array, 0, len(input))
	for _, v := range input {
		arr = append(arr, ValueOf(v))
	}
	return arr
}

func mapAsTable(input map[string]any) Table {
	t := make(Table, len(input))
	for k, v := range input {
		t[k] = ValueOf(v)
	}
	return t
}

func mapsAsTables(input []map[string]any) Tables {
	t := make(Tables, 0, len(input))
	for _, v := range input {
		t = append(t, mapAsTable(v))
	}
	return t
}

// --------------------------------------------------------------------

// lvalueOf converts the script input into a valid lua value
func lvalueOf(exec *lua.LState, value any) lua.LValue {
	if value == nil {
		return lua.LNil
	}

	switch val := reflect.ValueOf(value); val.Kind() {
	case reflect.Map, reflect.Slice:
		return ValueOf(val.Interface()).lvalue(exec)
	default:
		return luar.New(exec, value)
	}
}
