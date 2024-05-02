package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownPreference(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		clear(Options)
		newOpt("foo", "", "", OptionFlagPreference)

		p := Preferences{"foo": ""}
		assert.NoError(t, p.validate())
	})

	t.Run("existing but no preference", func(t *testing.T) {
		clear(Options)
		newOpt("foo", "", "", 0)

		p := Preferences{"foo": ""}
		assert.EqualError(t, p.validate(), "unknown preference: foo")
	})

	t.Run("not existing", func(t *testing.T) {
		clear(Options)
		p := Preferences{"foo": ""}
		assert.EqualError(t, p.validate(), "unknown preference: foo")
	})
}
