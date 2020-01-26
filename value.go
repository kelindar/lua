// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

// Value represents a returned
type Value interface {
	fmt.Stringer
	IsBoolean() bool
	IsNumber() bool
	IsString() bool
	ToBoolean() bool
	ToNumber() float64
	ToString() string
}

// ValueOf converts a value to a LUA-friendly one.
func ValueOf(i interface{}) Value {
	return value{luaValueOf(i)}
}

// value represents a LUA value
type value struct {
	lua.LValue
}

// String implements Stringer interface.
func (v value) String() string {
	return v.LValue.String()
}

// --------------------------------------------------------------------

// IsBoolean returns whether the value is a boolean.
func (v value) IsBoolean() bool {
	return v.LValue.Type() == lua.LTBool
}

// IsNumber returns whether the value is a number.
func (v value) IsNumber() bool {
	return v.LValue.Type() == lua.LTNumber
}

// IsString returns whether the value is a string.
func (v value) IsString() bool {
	return v.LValue.Type() == lua.LTString
}

// --------------------------------------------------------------------

// ToBoolean converts the value to a boolean value.
func (v value) ToBoolean() bool {
	return lua.LVAsBool(v.LValue)
}

// ToNumber converts the value to a numerical value.
func (v value) ToNumber() float64 {
	return float64(lua.LVAsNumber(v.LValue))
}

// ToString converts the value to a string value.
func (v value) ToString() string {
	return lua.LVAsString(v.LValue)
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
