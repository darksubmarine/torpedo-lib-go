package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlLoader_Load(t *testing.T) {
	m := Map{}

	loader := NewYamlLoader("_testing/config.yaml")
	loader.Load(m)

	assert.Equal(t, m.FetchStringP("log", "level"), "INFO")
	adapterListRaw, ok := m.Fetch("log", "adapters")
	assert.True(t, ok)

	adapterList, ok := adapterListRaw.([]interface{})
	assert.True(t, ok)
	assert.Len(t, adapterList, 2)

	adapter, ok := adapterList[0].(Map)
	assert.True(t, ok)
	assert.Equal(t, adapter.FetchStringP("type"), "stdout")
}
