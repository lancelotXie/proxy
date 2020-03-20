package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lancelot/proxy/proxy.lib/path"
)

func Test_SaveLoad(t *testing.T) {
	path.Init(path.Client)

	Set("123", 123)
	Set("3.2.1", 321)

	err := Save()
	if !assert.NoError(t, err) {
		return
	}

	Clear()
	err = Load()
	if !assert.NoError(t, err) {
		return
	}

	v0, ok := GetFloat("123")
	if !assert.True(t, ok) {
		return
	}
	if !assert.Equal(t, float64(123), v0) {
		return
	}

	v1, ok := GetInt("3.2.1")
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, 321, v1)
}
