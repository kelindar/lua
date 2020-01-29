// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

// This is a fork of https://github.com/layeh/gopher-json, licensed under The Unlicense

package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
)

func TestSimple(t *testing.T) {
	const str = `
	local json = require("json")
	assert(type(json) == "table")
	assert(type(json.decode) == "function")
	assert(type(json.encode) == "function")
	assert(json.encode(true) == "true")
	assert(json.encode(1) == "1")
	assert(json.encode(-10) == "-10")
	assert(json.encode(nil) == "null")
	assert(json.encode({}) == "[]")
	assert(json.encode({1, 2, 3}) == "[1,2,3]")
	local _, err = json.encode({1, 2, [10] = 3})
	assert(string.find(err, "sparse array"))
	local _, err = json.encode({1, 2, 3, name = "Tim"})
	assert(string.find(err, "mixed or invalid key types"))
	local _, err = json.encode({name = "Tim", [false] = 123})
	assert(string.find(err, "mixed or invalid key types"))
	local obj = {"a",1,"b",2,"c",3}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	for i = 1, #obj do
		assert(obj[i] == jsonObj[i])
	end
	local obj = {name="Tim",number=12345}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	assert(obj.name == jsonObj.name)
	assert(obj.number == jsonObj.number)
	assert(json.decode("null") == nil)
	assert(json.decode(json.encode({person={name = "tim",}})).person.name == "tim")
	local obj = {
		abc = 123,
		def = nil,
	}
	local obj2 = {
		obj = obj,
	}
	obj.obj2 = obj2
	assert(json.encode(obj) == nil)
	local a = {}
	for i=1, 5 do
		a[i] = i
	end
	assert(json.encode(a) == "[1,2,3,4,5]")
	`
	s := lua.NewState()
	defer s.Close()

	s.PreloadModule("json", Loader)
	assert.NotPanics(t, func() {
		s.DoString(str)
	})

}

func TestCustomRequire(t *testing.T) {
	const str = `
	local j = require("JSON")
	assert(type(j) == "table")
	assert(type(j.decode) == "function")
	assert(type(j.encode) == "function")
	`
	s := lua.NewState()
	defer s.Close()

	s.PreloadModule("JSON", Loader)
	if err := s.DoString(str); err != nil {
		t.Error(err)
	}
}

func TestDecodeValue_jsonNumber(t *testing.T) {
	s := lua.NewState()
	defer s.Close()

	v := DecodeValue(s, json.Number("124.11"))
	if v.Type() != lua.LTString || v.String() != "124.11" {
		t.Fatalf("expecting LString, got %T", v)
	}
}
