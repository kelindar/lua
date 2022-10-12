// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"encoding/json"
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

type numberType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

var (
	typeError   = reflect.TypeOf((*error)(nil)).Elem()
	typeNumber  = reflect.TypeOf(Number(0))
	typeString  = reflect.TypeOf(String(""))
	typeBool    = reflect.TypeOf(Bool(true))
	typeNumbers = reflect.TypeOf(Numbers(nil))
	typeStrings = reflect.TypeOf(Strings(nil))
	typeBools   = reflect.TypeOf(Bools(nil))
	typeTable   = reflect.TypeOf(Table(nil))
	typeTables  = reflect.TypeOf(Tables(nil))
)

var typeMap = map[reflect.Type]Type{
	typeString:  TypeString,
	typeNumber:  TypeNumber,
	typeBool:    TypeBool,
	typeStrings: TypeStrings,
	typeNumbers: TypeNumbers,
	typeBools:   TypeBools,
	typeTable:   TypeTable,
	typeTables:  TypeTables,
}

// Type represents a type of the value
type Type byte

// Various supported types
const (
	TypeNil = Type(iota)
	TypeBool
	TypeNumber
	TypeString
	TypeBools
	TypeNumbers
	TypeStrings
	TypeTable
	TypeTables
	TypeArray
)

// Value represents a returned
type Value interface {
	fmt.Stringer
	Type() Type
	Native() any
}

// ValueOf converts a value to a LUA-friendly one.
func ValueOf(i any) Value {
	return resultOf(luaValueOf(i))
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
				return resultOfTable(v)
			case lua.LTNumber:
				return resultOfNumbers(v)
			case lua.LTString:
				return resultOfStrings(v)
			case lua.LTBool:
				return resultOfBools(v)
			case lua.LTTable:
				return resultOfTables(v)
			}
		}
		return Nil{}
	case *lua.LUserData:
		return resultOfNative(v.Value)
	default:
		return Nil{}
	}
}

func resultOfNative(value any) Value {
	switch v := value.(type) {
	case map[string]any:
		return resultOfMap(v)
	case []map[string]any:
		return resultOfMaps(v)
	case []any:
		return resultOfArray(v)
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

		return resultOfNative(resp)
	}
}

func resultOfNumbers(t *lua.LTable) (out Numbers) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, float64(v.(lua.LNumber)))
	})
	return
}

func resultOfStrings(t *lua.LTable) (out Strings) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, string(v.(lua.LString)))
	})
	return
}

func resultOfBools(t *lua.LTable) (out Bools) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, bool(v.(lua.LBool)))
	})
	return
}

func resultOfTable(t *lua.LTable) Table {
	out := make(Table, t.Len())
	t.ForEach(func(k, v lua.LValue) {
		out[k.String()] = resultOf(v)
	})
	return out
}

func resultOfTables(t *lua.LTable) (out Tables) {
	t.ForEach(func(_, v lua.LValue) {
		out = append(out, resultOfTable(v.(*lua.LTable)))
	})
	return
}

func resultOfMap(input map[string]any) Table {
	t := make(Table, len(input))
	for k, v := range input {
		t[k] = ValueOf(v)
	}
	return t
}

func resultOfMaps(input []map[string]any) Tables {
	t := make(Tables, 0, len(input))
	for _, v := range input {
		t = append(t, resultOfMap(v))
	}
	return t
}

func resultOfArray(input []any) Tables {
	t := make(Tables, 0, len(input))
	for _, v := range input {
		t = append(t, resultOfMap(v.(map[string]any)))
	}
	return t
}

// --------------------------------------------------------------------

// luaValueOf converts a value to a LUA-friendly one.
func luaValueOf(i any) lua.LValue {
	switch v := i.(type) {
	case Number:
		return lua.LNumber(v)
	case String:
		return lua.LString(v)
	case Bool:
		return lua.LBool(v)
	case Numbers:
		return v.table()
	case Strings:
		return v.table()
	case Bools:
		return v.table()
	case Table:
		return v.table()
	case Tables:
		return v.table()
	case int:
		return lua.LNumber(v)
	case int8:
		return lua.LNumber(v)
	case int16:
		return lua.LNumber(v)
	case int32:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case uint:
		return lua.LNumber(v)
	case uint8:
		return lua.LNumber(v)
	case uint16:
		return lua.LNumber(v)
	case uint32:
		return lua.LNumber(v)
	case uint64:
		return lua.LNumber(v)
	case float32:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case map[string]any:
		return resultOfMap(v).table()
	case []int:
		return numbersOf(v).table()
	case []int8:
		return numbersOf(v).table()
	case []int16:
		return numbersOf(v).table()
	case []int32:
		return numbersOf(v).table()
	case []int64:
		return numbersOf(v).table()
	case []uint:
		return numbersOf(v).table()
	case []uint8:
		return numbersOf(v).table()
	case []uint16:
		return numbersOf(v).table()
	case []uint32:
		return numbersOf(v).table()
	case []uint64:
		return numbersOf(v).table()
	case []float32:
		return numbersOf(v).table()
	case []float64:
		return Numbers(v).table()
	case []bool:
		return Bools(v).table()
	case []string:
		return Strings(v).table()
	case []any:
		return luaArrayOf(v)
	case []map[string]any:
		return luaTablesOf(v)
	default:
		return lua.LNil
	}
}

// luaArrayOf creates a lua slice from an arbitrary array
func luaArrayOf(input []any) *lua.LTable {
	tbl := new(lua.LTable)
	for _, item := range input {
		tbl.Append(luaValueOf(item))
	}
	return tbl
}

// luaTablesOf creates a lua slice from a set of maps
func luaTablesOf(input []map[string]any) *lua.LTable {
	tbl := new(lua.LTable)
	for _, v := range input {
		tbl.Append(luaValueOf(v))
	}
	return tbl
}
