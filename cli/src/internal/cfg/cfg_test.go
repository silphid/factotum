package cfg

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/silphid/factotum/cli/src/internal/ctx"
	_assert "github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	assert := _assert.New(t)

	original := Config{
		Version: "version",
		Base: ctx.Context{
			Registry: ctx.RegistryECR,
			Image:    "image",
			Env: map[string]string{
				"ENV1": "value1",
				"ENV2": "value2",
			},
			Volumes: map[string]string{
				"/local/volume1": "/container/volume1",
				"/local/volume2": "/container/volume2",
			},
		},
		Contexts: map[string]ctx.Context{
			"ctx1": {
				Registry: ctx.RegistryECR,
				Image:    "image1",
				Env: map[string]string{
					"ENV1": "value1",
					"ENV2": "value2",
				},
				Volumes: map[string]string{
					"/local/volume1": "/container/volume1",
					"/local/volume2": "/container/volume2",
				},
			},
			"ctx2": {
				Registry: ctx.RegistryECR,
				Image:    "image1",
				Env: map[string]string{
					"ENV3": "value3",
					"ENV4": "value4",
				},
				Volumes: map[string]string{
					"/local/volume3": "/container/volume3",
					"/local/volume4": "/container/volume4",
				},
			},
		},
	}

	clone := original.Clone()

	diff := deep.Equal(original, clone)
	if diff != nil {
		assert.Fail("Cloned context different from original", diff)
	}

	assertNotSameMapStringString(t, original.Base.Env, clone.Base.Env)
	assertNotSameMapStringString(t, original.Base.Volumes, clone.Base.Volumes)
	assertNotSameMapStringContext(t, original.Contexts, clone.Contexts)
	assertNotSameMapStringString(t, original.Contexts["ctx1"].Env, clone.Contexts["ctx1"].Env)
	assertNotSameMapStringString(t, original.Contexts["ctx1"].Volumes, clone.Contexts["ctx1"].Volumes)
	assertNotSameMapStringString(t, original.Contexts["ctx2"].Env, clone.Contexts["ctx2"].Env)
	assertNotSameMapStringString(t, original.Contexts["ctx2"].Volumes, clone.Contexts["ctx2"].Volumes)
}

func TestMergingEmptyFieldsLeavesUnchanged(t *testing.T) {
	assert := _assert.New(t)

	config1 := Config{
		Version: "version1",
		Base: ctx.Context{
			Registry: ctx.RegistryECR,
			Image:    "image",
			Env: map[string]string{
				"ENV1": "value1",
				"ENV2": "value2",
			},
			Volumes: map[string]string{
				"/local/volume1": "/container/volume1",
				"/local/volume2": "/container/volume2",
			},
		},
		Contexts: map[string]ctx.Context{
			"ctx1": {
				Registry: ctx.RegistryGCR,
				Image:    "image1",
				Env: map[string]string{
					"ENV1": "value1",
					"ENV2": "value2",
				},
				Volumes: map[string]string{
					"/local/volume1": "/container/volume1",
					"/local/volume2": "/container/volume2",
				},
			},
			"ctx2": {
				Registry: ctx.RegistryDockerHub,
				Image:    "image2",
				Env: map[string]string{
					"ENV3": "value3",
					"ENV4": "value4",
				},
				Volumes: map[string]string{
					"/local/volume3": "/container/volume3",
					"/local/volume4": "/container/volume4",
				},
			},
		},
	}

	config2 := Config{}
	merged := config1.Merge(config2)

	diff := deep.Equal(merged, config1)
	if diff != nil {
		assert.Fail("Merged config different from expected config", diff)
	}
}

func TestMerge(t *testing.T) {
	assert := _assert.New(t)

	config1 := Config{
		Version: "version1",
		Base: ctx.Context{
			Registry: ctx.RegistryECR,
			Image:    "image",
			Env: map[string]string{
				"ENV1": "value1",
				"ENV2": "value2",
			},
			Volumes: map[string]string{
				"/local/volume1": "/container/volume1",
				"/local/volume2": "/container/volume2",
			},
		},
		Contexts: map[string]ctx.Context{
			"ctx1": {
				Registry: ctx.RegistryGCR,
				Image:    "image1",
				Env: map[string]string{
					"ENV1": "value1",
					"ENV2": "value2",
				},
				Volumes: map[string]string{
					"/local/volume1": "/container/volume1",
					"/local/volume2": "/container/volume2",
				},
			},
			"ctx2": {
				Registry: ctx.RegistryDockerHub,
				Image:    "image2",
				Env: map[string]string{
					"ENV3": "value3",
					"ENV4": "value4",
				},
				Volumes: map[string]string{
					"/local/volume3": "/container/volume3",
					"/local/volume4": "/container/volume4",
				},
			},
		},
	}

	config2 := Config{
		Version: "version2",
		Base: ctx.Context{
			Registry: ctx.RegistryGCR,
			Image:    "image2",
			Env: map[string]string{
				"ENV1": "value1b",
				"ENV3": "value3",
			},
			Volumes: map[string]string{
				"/local/volume1": "/container/volume1b",
				"/local/volume3": "/container/volume3",
			},
		},
		Contexts: map[string]ctx.Context{
			"ctx1": {
				Registry: ctx.RegistryECR,
				Image:    "image2",
				Env: map[string]string{
					"ENV1": "value1b",
					"ENV3": "value3",
				},
				Volumes: map[string]string{
					"/local/volume1": "/container/volume1b",
					"/local/volume3": "/container/volume3",
				},
			},
			"ctx3": {
				Registry: ctx.RegistryDockerHub,
				Image:    "image3",
				Env: map[string]string{
					"ENV5": "value5",
					"ENV6": "value6",
				},
				Volumes: map[string]string{
					"/local/volume5": "/container/volume5",
					"/local/volume6": "/container/volume6",
				},
			},
		},
	}

	expected := Config{
		Version: "version2",
		Base: ctx.Context{
			Registry: ctx.RegistryGCR,
			Image:    "image2",
			Env: map[string]string{
				"ENV1": "value1b",
				"ENV2": "value2",
				"ENV3": "value3",
			},
			Volumes: map[string]string{
				"/local/volume1": "/container/volume1b",
				"/local/volume2": "/container/volume2",
				"/local/volume3": "/container/volume3",
			},
		},
		Contexts: map[string]ctx.Context{
			"ctx1": {
				Registry: ctx.RegistryECR,
				Image:    "image2",
				Env: map[string]string{
					"ENV1": "value1b",
					"ENV2": "value2",
					"ENV3": "value3",
				},
				Volumes: map[string]string{
					"/local/volume1": "/container/volume1b",
					"/local/volume2": "/container/volume2",
					"/local/volume3": "/container/volume3",
				},
			},
			"ctx2": {
				Registry: ctx.RegistryDockerHub,
				Image:    "image2",
				Env: map[string]string{
					"ENV3": "value3",
					"ENV4": "value4",
				},
				Volumes: map[string]string{
					"/local/volume3": "/container/volume3",
					"/local/volume4": "/container/volume4",
				},
			},
			"ctx3": {
				Registry: ctx.RegistryDockerHub,
				Image:    "image3",
				Env: map[string]string{
					"ENV5": "value5",
					"ENV6": "value6",
				},
				Volumes: map[string]string{
					"/local/volume5": "/container/volume5",
					"/local/volume6": "/container/volume6",
				},
			},
		},
	}

	merged := config1.Merge(config2)

	diff := deep.Equal(merged, expected)
	if diff != nil {
		assert.Fail("Merged config different from expected config", diff)
	}
}

func assertNotSameMapStringString(t *testing.T, map1, map2 map[string]string, msgAndArgs ...interface{}) {
	_assert.NotEqual(t, reflect.ValueOf(map1).Pointer(), reflect.ValueOf(map2).Pointer(), msgAndArgs)
}

func assertNotSameMapStringContext(t *testing.T, map1, map2 map[string]ctx.Context, msgAndArgs ...interface{}) {
	_assert.NotEqual(t, reflect.ValueOf(map1).Pointer(), reflect.ValueOf(map2).Pointer(), msgAndArgs)
}
