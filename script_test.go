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
	assert.Equal(t, TypeNumber, out.Type())
	assert.Equal(t, Number(89), out.(Number))
	assert.Equal(t, "89", out.String())

	err = s.Close()
	assert.NoError(t, err)
}

// Benchmark_Serial/fib-8         	 6138716	       196 ns/op	      16 B/op	       2 allocs/op
// Benchmark_Serial/empty-8       	 7206128	       167 ns/op	       8 B/op	       1 allocs/op
func Benchmark_Serial(b *testing.B) {
	b.Run("fib", func(b *testing.B) {
		s, _ := newScript("fixtures/fib.lua")
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background(), 1)
		}
	})

	b.Run("empty", func(b *testing.B) {
		s, _ := newScript("fixtures/empty.lua")
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background())
		}
	})
}

// Benchmark_Fib_Parallel-8   	 4149006	       286 ns/op	      32 B/op	       3 allocs/op
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
