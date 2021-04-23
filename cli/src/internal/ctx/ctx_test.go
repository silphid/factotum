package ctx

import (
	"testing"

	"github.com/go-test/deep"
	_assert "github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	assert := _assert.New(t)

	context := Context{
		Env: map[string]string{
			"ENV1": "value1",
			"ENV2": "value2",
		},
		Volumes: map[string]string{
			"/local/volume1": "/container/volume1",
			"/local/volume2": "/container/volume2",
		},
	}

	clone := context.Clone()

	diff := deep.Equal(context, clone)
	if diff != nil {
		assert.Fail("Cloned context different from original", diff)
	}
}

func TestMerge(t *testing.T) {
	assert := _assert.New(t)

	context1 := Context{
		Env: map[string]string{
			"ENV1": "value1",
			"ENV2": "value2",
		},
		Volumes: map[string]string{
			"/local/volume1": "/container/volume1",
			"/local/volume2": "/container/volume2",
		},
	}

	context2 := Context{
		Env: map[string]string{
			"ENV1": "value1b",
			"ENV3": "value3",
		},
		Volumes: map[string]string{
			"/local/volume1": "/container/volume1b",
			"/local/volume3": "/container/volume3",
		},
	}

	expected := Context{
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
	}

	merged := context1.Merge(context2)

	diff := deep.Equal(merged, expected)
	if diff != nil {
		assert.Fail("Merged context different from expected context", diff)
	}
}
