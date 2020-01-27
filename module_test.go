// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"context"
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testModule() *Module {
	m := &Module{
		Name:    "test",
		Version: "1.0.0",
	}
	m.Register("hash", hash)
	m.Register("echo", echo)
	return m
}

func echo(v String) (String, error) {
	return v, nil
}

func hash(s String) (Number, error) {
	h := fnv.New32a()
	h.Write([]byte(s))

	return Number(h.Sum32()), nil
}

func Test_Hash(t *testing.T) {
	s, err := newScript("fixtures/hash.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), "abcdef")
	assert.NoError(t, err)
	assert.Equal(t, TypeNumber, out.Type())
	assert.Equal(t, int64(4282878506), int64(out.(Number)))

}

func Test_NotAFunc(t *testing.T) {
	m := &Module{
		Name:    "test",
		Version: "1.0.0",
	}
	assert.Error(t, m.Register("xxx", 123))
	assert.NoError(t, m.Register("hash", hash))
	assert.Equal(t, 1, len(m.funcs))
	m.Unregister("hash")
	assert.Equal(t, 0, len(m.funcs))
}
