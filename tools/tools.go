package tools

// This dedicated Go module is used to prevent bloating the root Go module with external
// tools.

import (
	_ "github.com/anchore/quill/cmd/quill"
	_ "github.com/jstemmer/go-junit-report/v2"
)
