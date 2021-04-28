package cfg

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/silphid/factotum/cli/src/internal/ctx"
	"github.com/silphid/factotum/cli/src/internal/helpers"
	_assert "github.com/stretchr/testify/assert"
)

func loadContext(file string) ctx.Context {
	path := filepath.Join("testdata", file+".yaml")
	if !helpers.PathExists(path) {
		panic(fmt.Errorf("context file not found: %q", path))
	}
	context, err := ctx.Load(path)
	if err != nil {
		panic(fmt.Errorf("loading context from %q: %w", path, err))
	}
	return context
}

func loadConfig(version, baseFile, ctx1Key, ctx1File, ctx2Key, ctx2File string) Config {
	return Config{
		Version: version,
		Base:    loadContext(baseFile),
		Contexts: map[string]ctx.Context{
			ctx1Key: loadContext(ctx1File),
			ctx2Key: loadContext(ctx2File),
		},
	}
}

func loadConfigWith3Contexts(version, baseFile, ctx1Key, ctx1File, ctx2Key, ctx2File, ctx3Key, ctx3File string) Config {
	return Config{
		Version: version,
		Base:    loadContext(baseFile),
		Contexts: map[string]ctx.Context{
			ctx1Key: loadContext(ctx1File),
			ctx2Key: loadContext(ctx2File),
			ctx3Key: loadContext(ctx3File),
		},
	}
}

func TestClone(t *testing.T) {
	assert := _assert.New(t)

	original := loadConfig("version1", "base1", "ctx1", "ctx1", "ctx2", "ctx2")
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

	config1 := loadConfig("version1", "base1", "ctx1", "ctx1", "ctx2", "ctx2")
	config2 := Config{}
	merged := config1.Merge(config2)

	diff := deep.Equal(merged, config1)
	if diff != nil {
		assert.Fail("Merged config different from expected config", diff)
	}
}

func TestMerge(t *testing.T) {
	assert := _assert.New(t)

	config1 := loadConfig("version1", "base1", "ctx1", "ctx1", "ctx2", "ctx2")
	config2 := loadConfig("version2", "base1b", "ctx1", "ctx1b", "ctx3", "ctx3")
	expected := loadConfigWith3Contexts("version2", "base1_base1b", "ctx1", "ctx1_ctx1b", "ctx2", "ctx2", "ctx3", "ctx3")
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

func TestResolveContextFromConfigs(t *testing.T) {
	assert := _assert.New(t)

	sharedConfig := loadConfig("version1", "base1", "ctx1", "ctx1", "ctx2", "ctx2")
	userConfig := loadConfig("version2", "base1b", "ctx1", "ctx1b", "ctx3", "ctx3")

	cases := []struct {
		name     string
		expected string
		error    string
	}{
		{
			name:     "ctx1",
			expected: "base1_base1b_ctx1_ctx1b",
		},
		{
			name:     "ctx2",
			expected: "base1_base1b_ctx2",
		},
		{
			name:     "ctx3",
			expected: "base1_base1b_ctx3",
		},
		{
			name:     "base",
			expected: "base1_base1b",
		},
		{
			name:     "none",
			expected: "none",
		},
		{
			name:  "not_found",
			error: "temp not found",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expected := loadContext(c.expected)
			actual, err := getContext(sharedConfig, userConfig, c.name)

			if c.error != "" {
				assert.NotNil(err)
				assert.Equal(c.error, err.Error())
			} else {
				assert.NoError(err)
				if diff := deep.Equal(expected, actual); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}
