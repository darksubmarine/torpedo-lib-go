package conf

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

var notExists = false
var exists = true

func TestConfMap_Fetch(t *testing.T) {
	m := Map{
		"L0-1": "value 1",
		"L0-2": 123,
		"L0-3": true,
		"L0-4": Map{
			"L1-1": "value level 1",
			"L1-2": 456,
			"L1-3": true,
			"L1-4": Map{
				"L2-1": "value level 2",
				"L2-2": 789,
				"L2-3": true,
				"L2-4": Map{
					"L3-1": uint8(123),
				},
			},
		},
	}

	hlpTestValue(t, m, notExists, nil, "level0")
	hlpTestValue(t, m, exists, "value 1", "L0-1")
	hlpTestValue(t, m, exists, 456, "L0-4", "L1-2")
	hlpTestValue(t, m, exists, uint8(123), "L0-4", "L1-4", "L2-4", "L3-1")

}

func hlpTestValue(t *testing.T, m Map, expected bool, val interface{}, key ...string) {
	v, ok := m.Fetch(key...)
	assert.Equal(t, val, v)
	assert.Equal(t, expected, ok)
}

func TestConfMap_Add(t *testing.T) {
	m := Map{
		"L0-1": "value 1",
		"L0-2": 123,
		"L0-3": true,
		"L0-4": Map{
			"L1-1": "value level 1",
			"L1-2": 456,
			"L1-3": true,
			"L1-4": Map{
				"L2-1": "value level 2",
				"L2-2": 789,
				"L2-3": true,
				"L2-4": Map{
					"L3-1": uint8(123),
				},
			},
		},
	}

	m.Add(789, "L0-1")
	m.Add("forgot bool", "L0-3")
	m.Add(float64(123.4567), "L00-4", "L1-4", "L2-4", "L3-2")
	fmt.Println(m.FetchP("L0-1"))

}
