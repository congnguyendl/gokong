package gokong

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

type Metal struct {
	Service *Id `json:"service" yaml:"service"`
}

func TestId_MarshalJSON(t *testing.T) {
	s := Metal{Service: ToId("123")}

	metal, err := json.Marshal(s)

	assert.Nil(t, err)
	assert.NotNil(t, metal)
	assert.EqualValues(t, metal, []byte(`{"service":{"id":"123"}}`))
}

func TestId_UnmarshalJSON(t *testing.T) {
	m := Metal{}

	err := json.Unmarshal([]byte(`{"service":{"id":"123"}}`), &m)

	assert.Nil(t, err)
	assert.EqualValues(t, m, Metal{Service: ToId("123")})
}

func TestId_MarshalYAML(t *testing.T) {
	s := Metal{Service: ToId("123")}

	metal, err := yaml.Marshal(s)

	assert.Nil(t, err)
	assert.NotNil(t, metal)
	assert.EqualValues(t, metal, []byte("service:\n  id: \"123\"\n"))
}

func TestId_UnmarshalYAML(t *testing.T) {
	m := Metal{}

	err := yaml.Unmarshal([]byte("service:\n  id: \"123\"\n"), &m)

	assert.Nil(t, err)
	assert.EqualValues(t, m, Metal{Service: ToId("123")})
}
