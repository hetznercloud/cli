//go:build tools
// +build tools

package tools

import (
	_ "github.com/anchore/quill/cmd/quill"
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/golang/mock/mockgen"
)
