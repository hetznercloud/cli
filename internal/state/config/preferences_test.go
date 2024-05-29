package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/state/config"
)

func TestPreferences_Get(t *testing.T) {
	t.Parallel()

	p := config.Preferences{
		"foo": "bar",
		"baz": "qux",
		"quux": map[string]any{
			"corge": "grault",
			"garply": map[string]any{
				"waldo": []string{"fred", "plugh"},
				"xyzzy": 2,
			},
		},
	}

	v, ok := p.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", v)

	v, ok = p.Get("baz")
	assert.True(t, ok)
	assert.Equal(t, "qux", v)

	v, ok = p.Get("buz")
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok = p.Get("quux.corge")
	assert.True(t, ok)
	assert.Equal(t, "grault", v)

	v, ok = p.Get("quux.garply.waldo")
	assert.True(t, ok)
	assert.Equal(t, []string{"fred", "plugh"}, v)

	v, ok = p.Get("quux.garply.xyzzy")
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	v, ok = p.Get("quux.garply")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{
		"waldo": []string{"fred", "plugh"},
		"xyzzy": 2,
	}, v)
}

func TestPreferences_Set(t *testing.T) {
	t.Parallel()

	p := config.Preferences{}
	p.Set("foo", "bar")
	p.Set("baz", "qux")
	p.Set("quux.corge", "grault")
	p.Set("quux.garply.waldo", []string{"fred", "plugh"})
	p.Set("quux.garply.xyzzy", 2)

	assert.Equal(t, config.Preferences{
		"foo": "bar",
		"baz": "qux",
		"quux": map[string]any{
			"corge": "grault",
			"garply": map[string]any{
				"waldo": []string{"fred", "plugh"},
				"xyzzy": 2,
			},
		},
	}, p)
}

func TestPreferences_Unset(t *testing.T) {
	t.Parallel()

	p := config.Preferences{
		"foo": "bar",
		"baz": "qux",
		"quux": map[string]any{
			"corge": "grault",
			"garply": map[string]any{
				"waldo": []string{"fred", "plugh"},
				"xyzzy": 2,
			},
		},
	}

	assert.True(t, p.Unset("foo"))
	assert.Equal(t, config.Preferences{
		"baz": "qux",
		"quux": map[string]any{
			"corge": "grault",
			"garply": map[string]any{
				"waldo": []string{"fred", "plugh"},
				"xyzzy": 2,
			},
		},
	}, p)

	assert.False(t, p.Unset("buz"))
	assert.Equal(t, config.Preferences{
		"baz": "qux",
		"quux": map[string]any{
			"corge": "grault",
			"garply": map[string]any{
				"waldo": []string{"fred", "plugh"},
				"xyzzy": 2,
			},
		},
	}, p)

	assert.True(t, p.Unset("quux.corge"))
	assert.Equal(t, config.Preferences{
		"baz": "qux",
		"quux": map[string]any{
			"garply": map[string]any{
				"waldo": []string{"fred", "plugh"},
				"xyzzy": 2,
			},
		},
	}, p)

	assert.True(t, p.Unset("quux.garply.waldo"))
	assert.Equal(t, config.Preferences{
		"baz": "qux",
		"quux": map[string]any{
			"garply": map[string]any{
				"xyzzy": 2,
			},
		},
	}, p)

	assert.True(t, p.Unset("quux.garply.xyzzy"))
	assert.Equal(t, config.Preferences{
		"baz": "qux",
	}, p)

	assert.True(t, p.Unset("baz"))
	assert.Equal(t, config.Preferences{}, p)
}

func TestPreferences_Validate(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		_, cleanup := config.NewTestOption("foo", "", "", config.OptionFlagPreference, nil)
		defer cleanup()

		p := config.Preferences{"foo": ""}
		assert.NoError(t, p.Validate())
	})

	t.Run("existing nested", func(t *testing.T) {
		_, cleanup := config.NewTestOption("foo.bar", "", "", config.OptionFlagPreference, nil)
		defer cleanup()

		p := config.Preferences{"foo": map[string]any{"bar": ""}}
		assert.NoError(t, p.Validate())
	})

	t.Run("existing but no preference", func(t *testing.T) {
		_, cleanup := config.NewTestOption("foo", "", "", 0, nil)
		defer cleanup()

		p := config.Preferences{"foo": ""}
		assert.EqualError(t, p.Validate(), "unknown preference: foo")
	})

	t.Run("existing nested but no preference", func(t *testing.T) {
		_, cleanup := config.NewTestOption("foo.bar", "", "", 0, nil)
		defer cleanup()

		p := config.Preferences{"foo": map[string]any{"bar": ""}}
		assert.EqualError(t, p.Validate(), "unknown preference: foo.bar")
	})

	t.Run("not existing", func(t *testing.T) {
		p := config.Preferences{"foo": ""}
		assert.EqualError(t, p.Validate(), "unknown preference: foo")
	})

	t.Run("not existing nested", func(t *testing.T) {
		p := config.Preferences{"foo": map[string]any{"bar": ""}}
		assert.EqualError(t, p.Validate(), "unknown preference: foo.bar")
	})

	t.Run("not existing deeply nested", func(t *testing.T) {
		_, cleanup := config.NewTestOption("foo.bar", "", "", 0, nil)
		defer cleanup()

		p := config.Preferences{"foo": map[string]any{"bar": map[string]any{"baz": ""}}}
		assert.EqualError(t, p.Validate(), "unknown preference: foo.bar.baz")
	})

	t.Run("nested missing map", func(t *testing.T) {
		_, cleanup := config.NewTestOption("foo.bar.baz", "", "", 0, nil)
		defer cleanup()

		p := config.Preferences{"foo": map[string]any{"bar": "this should be a map"}}
		assert.EqualError(t, p.Validate(), "unknown preference: foo.bar")
	})
}
