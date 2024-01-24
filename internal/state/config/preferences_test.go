package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownPreference(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		clear(opts)
		newOpt("foo", "", "", OptionSourcePreference)

		p := preferences{"foo": ""}
		assert.NoError(t, p.validate())
	})

	t.Run("existing but no preference", func(t *testing.T) {
		clear(opts)
		newOpt("foo", "", "", 0)

		p := preferences{"foo": ""}
		assert.EqualError(t, p.validate(), "unknown preference: foo")
	})

	t.Run("not existing", func(t *testing.T) {
		clear(opts)
		p := preferences{"foo": ""}
		assert.EqualError(t, p.validate(), "unknown preference: foo")
	})
}
