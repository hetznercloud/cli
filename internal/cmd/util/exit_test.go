package util_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
)

func TestExitError(t *testing.T) {
	cause := errors.New("process failed")
	err := util.NewExitError(23, cause, true)

	assert.Equal(t, 23, util.ExitCode(err))
	assert.True(t, util.IsSilent(err))
	assert.ErrorIs(t, err, cause)
}

func TestExitErrorDefaultsInvalidCode(t *testing.T) {
	err := util.NewExitError(-1, errors.New("process failed"), false)

	assert.Equal(t, 1, util.ExitCode(err))
	assert.False(t, util.IsSilent(err))
}
