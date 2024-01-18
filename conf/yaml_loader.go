package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

type YamlLoader struct {
	filename string
}

func NewYamlLoader(filename string) *YamlLoader {
	return &YamlLoader{filename: filename}
}

func (y *YamlLoader) Load(m Map) Map {
	if data, err := os.ReadFile(y.filename); err != nil {
		panic(err)
	} else {
		if err := yaml.Unmarshal(data, &m); err != nil {
			panic(err)
		}
	}

	return m
}
