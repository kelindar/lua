// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package lua

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestValue(typ Type) Value {
	switch typ {
	case TypeNumber:
		return Number(1)
	case TypeString:
		return String("x")
	case TypeBool:
		return Bool(true)
	default:
		return Nil{}
	}
}

func TestTableConvert(t *testing.T) {
	expect := newTestMap()
	input := ValueOf(expect)
	assert.EqualValues(t, expect, input.Native())
}

func TestTableCodec(t *testing.T) {
	expect := newTestMap()
	encoded, err := json.Marshal(expect)
	assert.NoError(t, err)

	// Decode a table
	var decoded Table
	assert.NoError(t, json.Unmarshal(encoded, &decoded))

	// Reincode it back to json and compare strings
	reincoded, err := json.Marshal(decoded)
	assert.NoError(t, err)
	assert.Equal(t, encoded, reincoded)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func newTestMap() map[string]any {
	return map[string]any{
		"user":   "Roman",
		"age":    37.0,
		"dev":    true,
		"bitmap": []bool{true, false, true},
		"floats": []float64{1, 2, 3, 4, 5},
		"skills": map[string]any{
			"golang": 52.7,
			"eating": 100.0,
			"tables": []map[string]any{{
				"table": 1.0,
				"args":  []float64{2, 4},
			}, {
				"table": 2.0,
				"args":  []float64{2, 4},
			}},
		},
	}
}
