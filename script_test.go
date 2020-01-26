// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newScript(file string) (*Script, error) {
	f, _ := os.Open(file)
	return FromReader("test.lua", f)
}

func TestGet(t *testing.T) {

	s, err := newScript("fixtures/fib.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), 10)
	assert.NoError(t, err)
	assert.True(t, out.IsNumber())
	assert.Equal(t, float64(89), out.ToNumber())
	assert.Equal(t, "89", out.String())

	err = s.Close()
	assert.NoError(t, err)
}

func Benchmark_Fib_Serial(b *testing.B) {
	s, _ := newScript("fixtures/fib.lua")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.Run(context.Background(), 1)
	}
}

func Benchmark_Fib_Parallel(b *testing.B) {
	s, _ := newScript("fixtures/fib.lua")
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Run(context.Background(), 1)
		}
	})
}
