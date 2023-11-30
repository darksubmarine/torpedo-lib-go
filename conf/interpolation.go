package conf

import (
	"os"
	"regexp"
)

const (
	envIdxName         = 1
	envIdxDefaultValue = 2
)

var yamlEnvRegex = regexp.MustCompile(`^\s*\$\{(\w+)(:\w+)?\}\s*$`)

func processWithEnvVarInterpolation(raw string) string {

	match := yamlEnvRegex.FindAllStringSubmatch(raw, 1)
	if len(match) == 0 || len(match[0]) < 3 { // no match or invalid match/emtpy capturing groups
		return raw
	}

	groups := match[0]
	name := groups[envIdxName]
	fallback := groups[envIdxDefaultValue]
	value := os.Getenv(name)

	// if the value of the env var is empty, but there is a fallback, return it
	if len(value) == 0 && len(fallback) > 1 {
		return fallback[1:]
	}

	return value
}
