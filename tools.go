//go:build tools
// +build tools

package tools

// When updating this, think about also updating script/install_tools.sh
import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/go-cmp/cmp"
	_ "github.com/rjeczalik/interfaces/cmd/interfacer"
)
