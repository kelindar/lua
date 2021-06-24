// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newScript(file string, mods ...Module) (*Script, error) {
	f, _ := os.Open(file)
	mods = append(mods, testModule())
	return FromReader("test.lua", f, mods...)
}

type Person struct {
	Name string
	Age  int
}

// Benchmark_Serial/fib-8         	 5765884	       207 ns/op	      16 B/op	       2 allocs/op
// Benchmark_Serial/empty-8       	 8471679	       142 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Serial/update-8      	 1000000	      1140 ns/op	     256 B/op	      14 allocs/op
// Benchmark_Serial/sleep-sigle-8 	     100	  10355408 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Serial/sleep-multi-8 	     890	   1349583 ns/op	       0 B/op	       0 allocs/op
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

	b.Run("update", func(b *testing.B) {
		s, _ := newScript("fixtures/update.lua")
		input := &Person{Name: "Roman"}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background(), input)
		}
	})

	b.Run("sleep-sigle", func(b *testing.B) {
		s, _ := newScript("fixtures/sleep.lua")
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background())
		}
	})

	b.Run("sleep-multi", func(b *testing.B) {
		s, _ := newScript("fixtures/sleep.lua")
		b.RunParallel(func(pb *testing.PB) {
			b.ReportAllocs()
			b.ResetTimer()
			for pb.Next() {
				s.Run(context.Background())
			}
		})
	})

}

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
Benchmark_Module/echo-12         	 2925387	       405.3 ns/op	      48 B/op	       3 allocs/op
Benchmark_Module/hash-12         	 3037982	       388.2 ns/op	      32 B/op	       3 allocs/op
*/
func Benchmark_Module(b *testing.B) {
	b.Run("echo", func(b *testing.B) {
		s, _ := newScript("fixtures/echo.lua")
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background(), "abc")
		}
	})

	b.Run("hash", func(b *testing.B) {
		s, _ := newScript("fixtures/hash.lua")
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background(), "abc")
		}
	})
}

// Benchmark_Fib_Parallel-8   	 3893732	       268 ns/op	      16 B/op	       2 allocs/op
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

func Test_Update(t *testing.T) {
	s, err := newScript("fixtures/update.lua")
	assert.NoError(t, err)
	assert.Equal(t, "test.lua", s.Name())

	input := &Person{
		Name: "Roman",
	}
	out, err := s.Run(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, TypeString, out.Type())
	assert.Equal(t, "Updated", input.Name)
	assert.Equal(t, "Updated", out.String())
}

func Test_Fib(t *testing.T) {

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

func Test_Empty(t *testing.T) {
	s, err := newScript("fixtures/empty.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, TypeNil, out.Type())
}

func Test_Print(t *testing.T) {
	s, err := newScript("fixtures/print.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), &Person{
		Name: "Roman",
	})
	assert.NoError(t, err)
	assert.Equal(t, TypeString, out.Type())
	assert.Equal(t, "Hello, Roman!", out.String())
}

func Test_InvalidScript(t *testing.T) {
	_, err := FromString("", `
	xxx main()
		local x = 1
	end`)
	assert.Error(t, err)
}

func Test_NoMain(t *testing.T) {

	{
		s, err := FromString("", `main = 1`)
		assert.NoError(t, err)

		_, err = s.Run(context.Background())
		assert.Error(t, err)
	}

	{
		s, err := FromString("", `
		function notmain()
			local x = 1
		end`)
		assert.NoError(t, err)

		_, err = s.Run(context.Background())
		assert.Error(t, err)
	}

	{
		s, err := FromString("", `
		function xxx()
			local x = 1
		end`)
		assert.NoError(t, err)

		_, err = s.Run(context.Background())
		assert.Error(t, err)
	}
}

func Test_Error(t *testing.T) {
	{
		_, err := FromString("", `
		error() 
		function main()
			local x = 1
		end`)
		assert.Error(t, err)
	}

	{
		s, err := FromString("", `
		function main()
			error() 
		end`)
		assert.NoError(t, err)
		_, err = s.Run(context.Background())
		assert.Error(t, err)
	}
}

func Test_JSON(t *testing.T) {
	input := map[string]interface{}{
		"a": 123,
		"b": "hello",
		"c": 10.15,
		"d": true,
		"e": &Person{Name: "Roman", Age: 15},
	}

	s, err := newScript("fixtures/json.lua")
	assert.NoError(t, err)

	out, err := s.Run(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, TypeString, out.Type())
	assert.Equal(t, `{"a":123,"b":"hello","c":10.15,"d":true,"e":{"Name":"Roman","Age":15}}`,
		out.String())
}
