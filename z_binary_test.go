// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package lua

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringString(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v String) (String, error) {
		return newTestValue(TypeString).(String), nil
	})
	m.Register("test2", func(v String) (String, error) {
		return newTestValue(TypeString).(String), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeString).(String))
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
		_, err = s.Run(context.Background(), newTestValue(TypeString).(String))
		assert.Error(t, err)
	}
}

func Test_StringNumber(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v String) (Number, error) {
		return newTestValue(TypeNumber).(Number), nil
	})
	m.Register("test2", func(v String) (Number, error) {
		return newTestValue(TypeNumber).(Number), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeString).(String))
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
		_, err = s.Run(context.Background(), newTestValue(TypeString).(String))
		assert.Error(t, err)
	}
}

func Test_StringBool(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v String) (Bool, error) {
		return newTestValue(TypeBool).(Bool), nil
	})
	m.Register("test2", func(v String) (Bool, error) {
		return newTestValue(TypeBool).(Bool), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeString).(String))
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
		_, err = s.Run(context.Background(), newTestValue(TypeString).(String))
		assert.Error(t, err)
	}
}

func Test_NumberString(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v Number) (String, error) {
		return newTestValue(TypeString).(String), nil
	})
	m.Register("test2", func(v Number) (String, error) {
		return newTestValue(TypeString).(String), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeNumber).(Number))
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
		_, err = s.Run(context.Background(), newTestValue(TypeNumber).(Number))
		assert.Error(t, err)
	}
}

func Test_NumberNumber(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v Number) (Number, error) {
		return newTestValue(TypeNumber).(Number), nil
	})
	m.Register("test2", func(v Number) (Number, error) {
		return newTestValue(TypeNumber).(Number), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeNumber).(Number))
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
		_, err = s.Run(context.Background(), newTestValue(TypeNumber).(Number))
		assert.Error(t, err)
	}
}

func Test_NumberBool(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v Number) (Bool, error) {
		return newTestValue(TypeBool).(Bool), nil
	})
	m.Register("test2", func(v Number) (Bool, error) {
		return newTestValue(TypeBool).(Bool), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeNumber).(Number))
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
		_, err = s.Run(context.Background(), newTestValue(TypeNumber).(Number))
		assert.Error(t, err)
	}
}

func Test_BoolString(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v Bool) (String, error) {
		return newTestValue(TypeString).(String), nil
	})
	m.Register("test2", func(v Bool) (String, error) {
		return newTestValue(TypeString).(String), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeBool).(Bool))
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
		_, err = s.Run(context.Background(), newTestValue(TypeBool).(Bool))
		assert.Error(t, err)
	}
}

func Test_BoolNumber(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v Bool) (Number, error) {
		return newTestValue(TypeNumber).(Number), nil
	})
	m.Register("test2", func(v Bool) (Number, error) {
		return newTestValue(TypeNumber).(Number), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeBool).(Bool))
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
		_, err = s.Run(context.Background(), newTestValue(TypeBool).(Bool))
		assert.Error(t, err)
	}
}

func Test_BoolBool(t *testing.T) {
	m := &NativeModule{Name: "test"}
	m.Register("test1", func(v Bool) (Bool, error) {
		return newTestValue(TypeBool).(Bool), nil
	})
	m.Register("test2", func(v Bool) (Bool, error) {
		return newTestValue(TypeBool).(Bool), errors.New("boom")
	})

	{ // Happy path
		s, err := FromString("", `
		local api = require("test")
		function main(input)
			return api.test1(input)
		end`, m)
		assert.NotNil(t, s)
		assert.NoError(t, err)
		_, err = s.Run(context.Background(), newTestValue(TypeBool).(Bool))
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
		_, err = s.Run(context.Background(), newTestValue(TypeBool).(Bool))
		assert.Error(t, err)
	}
}
