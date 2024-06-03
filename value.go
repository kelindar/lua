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

type tableType interface {
	lcopy(*lua.LTable, *lua.LState) lua.LValue
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
	typeArray   = reflect.TypeOf(Array(nil))
	typeValue   = reflect.TypeOf((*Value)(nil)).Elem()
)

var typeMap = map[reflect.Type]Type{
	typeString:  TypeString,
	typeNumber:  TypeNumber,
	typeBool:    TypeBool,
	typeStrings: TypeStrings,
	typeNumbers: TypeNumbers,
	typeBools:   TypeBools,
	typeTable:   TypeTable,
	typeArray:   TypeArray,
	typeValue:   TypeValue,
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
	TypeArray
	TypeValue
)

// Value represents a returned
type Value interface {
	fmt.Stringer
	Type() Type
	Native() any
	lvalue(*lua.LState) lua.LValue
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

// lvalue converts the value to a LUA value
func (v Nil) lvalue(*lua.LState) lua.LValue {
	return lua.LNil
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

// lvalue converts the value to a LUA value
func (v Number) lvalue(*lua.LState) lua.LValue {
	return lua.LNumber(v)
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

// lvalue converts the value to a LUA value
func (v Numbers) lvalue(state *lua.LState) lua.LValue {
	tbl := state.CreateTable(len(v)+4, 0)
	return v.lcopy(tbl, state)
}

// lcopy copies the table to another table
func (v Numbers) lcopy(dst *lua.LTable, state *lua.LState) lua.LValue {
	for _, item := range v {
		dst.Append(lua.LNumber(item))
	}
	return dst
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

// MarshalJSON marshals string to byte
func (v String) MarshalJSON() ([]byte, error) {
	// if it's below special string we marshal as empty array
	if v == "[e7d47667-b92a-48b5-977a-b3199ab09ff9]" {
		return json.Marshal([]any{})
	}
	return json.Marshal(string(v))
}

// lvalue converts the value to a LUA value
func (v String) lvalue(*lua.LState) lua.LValue {
	return lua.LString(v)
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

// lvalue converts the value to a LUA value
func (v Strings) lvalue(state *lua.LState) lua.LValue {
	tbl := state.CreateTable(len(v)+4, 0)
	return v.lcopy(tbl, state)
}

// lcopy copies the table to another table
func (v Strings) lcopy(dst *lua.LTable, state *lua.LState) lua.LValue {
	for _, item := range v {
		dst.Append(lua.LString(item))
	}
	return dst
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

// lvalue converts the value to a LUA value
func (v Bool) lvalue(*lua.LState) lua.LValue {
	return lua.LBool(v)
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

// lvalue converts the value to a LUA value
func (v Bools) lvalue(state *lua.LState) lua.LValue {
	tbl := state.CreateTable(len(v)+4, 0)
	return v.lcopy(tbl, state)
}

// lcopy copies the table to another table
func (v Bools) lcopy(dst *lua.LTable, state *lua.LState) lua.LValue {
	for _, item := range v {
		dst.Append(lua.LBool(item))
	}
	return dst
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
	out := make(map[string]any, len(v))
	for key, elem := range v {
		out[key] = elem.Native()
	}
	return out
}

// lvalue converts the value to a LUA value
func (v Table) lvalue(state *lua.LState) lua.LValue {
	tbl := state.CreateTable(0, len(v)+4)
	return v.lcopy(tbl, state)
}

// lcopy copies the table to another table
func (v Table) lcopy(dst *lua.LTable, state *lua.LState) lua.LValue {
	for k, item := range v {
		dst.RawSetString(k, item.lvalue(state))
	}
	return dst
}

// UnmarshalJSON unmarshals the type from JSON
func (v *Table) UnmarshalJSON(b []byte) error {
	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*v = mapAsTable(data)
	return nil
}

// --------------------------------------------------------------------

// Array represents the array of values
type Array []Value

// Type returns the type of the value
func (v Array) Type() Type {
	return TypeArray
}

// String returns the string representation of the value
func (v Array) String() string {
	return fmt.Sprintf("%+v", []Value(v))
}

// Native returns value casted to native type
func (v Array) Native() any {
	out := make([]any, len(v))
	for i, elem := range v {
		out[i] = elem.Native()
	}
	return out
}

// lvalue converts the value to a LUA value
func (v Array) lvalue(state *lua.LState) lua.LValue {
	tbl := state.CreateTable(len(v)+4, 0)
	return v.lcopy(tbl, state)
}

// lcopy copies the table to another table
func (v Array) lcopy(dst *lua.LTable, state *lua.LState) lua.LValue {
	for _, item := range v {
		dst.Append(item.lvalue(state))
	}
	return dst
}
