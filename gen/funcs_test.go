package lua

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate genny -in=$GOFILE -out=../z_funcs_test.go gen "TIn=String,Number,Bool TOut=String,Number,Bool"

func funcTInTOut(v TIn) (TOut, error) {
	return newTestValue(TypeTOut).(TOut), nil
}

func funcTInTOutWithErr(v TIn) (TOut, error) {
	return newTestValue(TypeTOut).(TOut), errors.New("boom")
}

func Test_TInTOut(t *testing.T) {
	m := &Module{Name: "test"}
	m.Register("test1", funcTInTOut)
	m.Register("test2", funcTInTOutWithErr)

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
