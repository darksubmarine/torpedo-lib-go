package conf

import (
	"testing"
)

type TestVal struct {
	expected string
	raw      string
}

var toTest = []TestVal{
	{expected: "http://localhost:8090/api/v1/some/endpoint?q=1&q#=3", raw: "${ENV_LOCATION:http://localhost:8090/api/v1/some/endpoint?q=1&q#=3}"},
	{expected: "someValue98-_", raw: "${ENV_LOCATION:someValue98-_}"},
	{expected: "some@email.com", raw: "${ENV_LOCATION:some@email.com}"},
	{expected: "invalid", raw: "${ENV_LOCATION:invalid}"},
	{expected: "", raw: "${ENV_LOCATION}"},
}

func TestMap_Interpolate(t *testing.T) {
	for _, test := range toTest {
		got := processWithEnvVarInterpolation(test.raw)
		if test.expected != got {
			t.Errorf("expected: %s\nactual: %s", test.expected, got)
		}
	}
}
