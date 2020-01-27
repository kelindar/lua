// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"fmt"
	"reflect"

	"github.com/yuin/gopher-lua"
)

var (
	typeError  = reflect.TypeOf((*error)(nil)).Elem()
	typeValue  = reflect.TypeOf((*Value)(nil)).Elem()
	typeNumber = reflect.TypeOf(Number(0))
	typeString = reflect.TypeOf(String(""))
	typeBool   = reflect.TypeOf(Bool(true))
)

var typeMap = map[reflect.Type]Type{
	typeString: TypeString,
	typeNumber: TypeNumber,
	typeBool:   TypeBool,
}

// Type represents a type of the value
type Type byte

// Various supported types
const (
	TypeNil = Type(iota)
	TypeBool
	TypeNumber
	TypeString
)

// Value represents a returned
type Value interface {
	fmt.Stringer
	Type() Type
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
	default:
		return Nil{}
	}
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
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case reflect.Value:
		switch v.Type() {
		case typeString:
			return lua.LString(v.String())
		case typeNumber:
			return lua.LNumber(v.Float())
		case typeBool:
			return lua.LBool(v.Bool())
		default:
			return lua.LNil
		}
	default:
		return lua.LNil
	}
}
