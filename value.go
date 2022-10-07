// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
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
)

var typeMap = map[reflect.Type]Type{
	typeString:  TypeString,
	typeNumber:  TypeNumber,
	typeBool:    TypeBool,
	typeStrings: TypeStrings,
	typeNumbers: TypeNumbers,
	typeBools:   TypeBools,
	typeTable:   TypeTable,
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
)

// Value represents a returned
type Value interface {
	fmt.Stringer
	Type() Type
	Native() interface{}
}

// ValueOf converts a value to a LUA-friendly one.
func ValueOf(i interface{}) Value {
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
		out[string(k.String())] = resultOf(v)
	})
	return out
}

func resultOfMap(input map[string]interface{}) Table {
	t := make(Table, len(input))
	for k, v := range input {
		t[k] = ValueOf(v)
	}
	return t
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
func (v Number) Native() interface{} {
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
func (v Numbers) Native() interface{} {
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
func (v String) Native() interface{} {
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
func (v Strings) Native() interface{} {
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
func (v Bool) Native() interface{} {
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
func (v Bools) Native() interface{} {
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
func (v Table) Native() interface{} {
	out := make(map[string]interface{})
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
func (v Nil) Native() interface{} {
	return nil
}

// --------------------------------------------------------------------

// luaValueOf converts a value to a LUA-friendly one.
func luaValueOf(i interface{}) lua.LValue {
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
	case map[string]interface{}:
		return resultOfMap(v).table()
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
		default:
			return lua.LNil
		}
	default:
		return lua.LNil
	}
}
