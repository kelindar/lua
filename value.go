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
	switch v.Type() {
	case lua.LTNumber:
		return Number(v.(lua.LNumber))
	case lua.LTString:
		return String(v.(lua.LString))
	case lua.LTBool:
		return Bool(v.(lua.LBool))
	case lua.LTTable:
		table := v.(*lua.LTable)
		if top := table.RawGetInt(1); top != nil {
			switch top.Type() {
			case lua.LTNil:
				return resultOfTable(table)
			case lua.LTNumber:
				return resultOfNumbers(table)
			case lua.LTString:
				return resultOfStrings(table)
			case lua.LTBool:
				return resultOfBools(table)
			case lua.LTTable:
				return resultOfTables(table)
			}
		}

		return Nil{}
	default:
		return Nil{}
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

// --------------------------------------------------------------------

// Nil represents the nil value
type Nil struct{}

// Type returns the type of the value
func (v Nil) Type() Type {
	return TypeNil
}

// String returns the string representation of the value
func (v Nil) String() string {
	return "(nil)"
}

// Native returns value casted to native type
func (v Nil) Native() any {
	return nil
}

// --------------------------------------------------------------------

// Number represents the numerical value
type Number float64

// Type returns the type of the value
func (v Number) Type() Type {
	return TypeNumber
}

// String returns the string representation of the value
func (v Number) String() string {
	return lua.LNumber(v).String()
}

// Native returns value casted to native type
func (v Number) Native() any {
	return float64(v)
}

// --------------------------------------------------------------------

// Numbers represents the number array value
type Numbers []float64

// numbersOf returns an array as a numbers array
func numbersOf[T numberType](arr []T) Numbers {
	out := make([]float64, 0, len(arr))
	for _, v := range arr {
		out = append(out, float64(v))
	}
	return out
}

// Type returns the type of the value
func (v Numbers) Type() Type {
	return TypeNumbers
}

// String returns the string representation of the value
func (v Numbers) String() string {
	return fmt.Sprintf("%v", []float64(v))
}

// Native returns value casted to native type
func (v Numbers) Native() any {
	return []float64(v)
}

// Table converts the slice to a lua table
func (v Numbers) table() *lua.LTable {
	tbl := new(lua.LTable)
	for _, item := range v {
		tbl.Append(lua.LNumber(item))
	}
	return tbl
}

// --------------------------------------------------------------------

// String represents the string value
type String string

// Type returns the type of the value
func (v String) Type() Type {
	return TypeString
}

// String returns the string representation of the value
func (v String) String() string {
	return lua.LString(v).String()
}

// Native returns value casted to native type
func (v String) Native() any {
	return string(v)
}

// --------------------------------------------------------------------

// Strings represents the string array value
type Strings []string

// Type returns the type of the value
func (v Strings) Type() Type {
	return TypeStrings
}

// String returns the string representation of the value
func (v Strings) String() string {
	return fmt.Sprintf("%v", []string(v))
}

// Native returns value casted to native type
func (v Strings) Native() any {
	return []string(v)
}

// Table converts the slice to a lua table
func (v Strings) table() *lua.LTable {
	tbl := new(lua.LTable)
	for _, item := range v {
		tbl.Append(lua.LString(item))
	}
	return tbl
}

// --------------------------------------------------------------------

// Bool represents the boolean value
type Bool bool

// Type returns the type of the value
func (v Bool) Type() Type {
	return TypeBool
}

// String returns the string representation of the value
func (v Bool) String() string {
	return lua.LBool(v).String()
}

// Native returns value casted to native type
func (v Bool) Native() any {
	return bool(v)
}

// --------------------------------------------------------------------

// Bools represents the boolean array value
type Bools []bool

// Type returns the type of the value
func (v Bools) Type() Type {
	return TypeBools
}

// String returns the string representation of the value
func (v Bools) String() string {
	return fmt.Sprintf("%v", []bool(v))
}

// Native returns value casted to native type
func (v Bools) Native() any {
	return []bool(v)
}

// Table converts the slice to a lua table
func (v Bools) table() *lua.LTable {
	tbl := new(lua.LTable)
	for _, item := range v {
		tbl.Append(lua.LBool(item))
	}
	return tbl
}

// --------------------------------------------------------------------

// Table represents a map of string to value
type Table map[string]Value

// Type returns the type of the value
func (v Table) Type() Type {
	return TypeTable
}

// String returns the string representation of the value
func (v Table) String() string {
	return fmt.Sprintf("%v", map[string]Value(v))
}

// Native returns value casted to native type
func (v Table) Native() any {
	out := make(map[string]any)
	for key, elem := range v {
		out[key] = elem.Native()
	}
	return out
}

// Table converts the slice to a lua table
func (v Table) table() *lua.LTable {
	tbl := new(lua.LTable)
	for k, item := range v {
		tbl.RawSetString(k, luaValueOf(item))
	}
	return tbl
}

// UnmarshalJSON unmarshals the type from JSON
func (v *Table) UnmarshalJSON(b []byte) error {
	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*v = resultOfMap(data)
	return nil
}

// --------------------------------------------------------------------

// Tables represents the array of tables
type Tables []Table

// Type returns the type of the value
func (v Tables) Type() Type {
	return TypeTables
}

// String returns the string representation of the value
func (v Tables) String() string {
	return fmt.Sprintf("%v", "(array of tables)")
}

// Native returns value casted to native type
func (v Tables) Native() any {
	var out []map[string]any
	for _, elem := range v {
		if tbl, ok := elem.Native().(map[string]any); ok {
			out = append(out, tbl)
		}
	}
	return out
}

// Table converts the slice to a lua table
func (v Tables) table() *lua.LTable {
	tbl := new(lua.LTable)
	for _, item := range v {
		tbl.Append(luaValueOf(item))
	}
	return tbl
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
	case reflect.Value:
		switch v.Type() {
		case typeString:
			return lua.LString(v.String())
		case typeNumber:
			return lua.LNumber(v.Float())
		case typeBool:
			return lua.LBool(v.Bool())
		case typeNumbers:
			return v.Interface().(Numbers).table()
		case typeBools:
			return v.Interface().(Bools).table()
		case typeStrings:
			return v.Interface().(Strings).table()
		case typeTable:
			return v.Interface().(Table).table()
		case typeTables:
			return v.Interface().(Tables).table()
		default:
			return lua.LNil
		}

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
