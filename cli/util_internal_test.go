package cli

import "testing"

func TestOnlyOneSet(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		ss       []string
		expected bool
	}{
		{
			name:     "only arg emtpy",
			expected: false,
		},
		{
			name:     "only arg non-empty",
			s:        "s",
			expected: true,
		},
		{
			name:     "first arg non-empty, rest empty",
			s:        "s",
			ss:       []string{""},
			expected: true,
		},
		{
			name: "at least one other arg non-empty",
			s:    "s",
			ss:   []string{"", "s"},
		},
		{
			name:     "only one arg non-empty",
			ss:       []string{"", "s"},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := exactlyOneSet(tt.s, tt.ss...)
			if tt.expected != actual {
				t.Errorf("expected %t; got %t", tt.expected, actual)
			}
		})
	}
}
