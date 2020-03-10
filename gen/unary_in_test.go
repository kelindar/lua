package lua

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate genny -in=$GOFILE -out=../z_unary_test.go gen "TIn=String,Number,Bool"

func Test_In_TIn(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v TIn) error {
		return nil
	})
	m.Register("test2", func(v TIn) error {
		return errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeTIn).(TIn))
		assert.NoError(t, err)
	}

	{ // Invalid argument
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test2(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeTIn).(TIn))
		assert.Error(t, err)
	}
}

func Test_Out_TIn(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func() (TIn, error) {
		return newTestValue(TypeTIn).(TIn), nil
	})
	m.Register("test2", func() (TIn, error) {
		return newTestValue(TypeTIn).(TIn), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background())
		assert.NoError(t, err)
	}

	{ // Invalid argument
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test2(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background())
		assert.Error(t, err)
	}
}
