// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

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
func resultOf(i lua.LValue) Value {
	switch v := i.(type) {
	case lua.LNumber:
		return Number(v)
	case lua.LString:
		return String(v)
	case lua.LBool:
		return Bool(v)
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

// luaValuesOf converts a set of values to a LUA-friendly ones.
func luaValuesOf(i []interface{}) []lua.LValue {
	out := make([]lua.LValue, 0, len(i))
	for _, v := range i {
		out = append(out, luaValueOf(v))
	}
	return out
}

// luaValueOf converts a value to a LUA-friendly one.
func luaValueOf(i interface{}) lua.LValue {
	switch v := i.(type) {
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
	default:
		return lua.LNil
	}
}
