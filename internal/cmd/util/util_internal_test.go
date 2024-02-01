package util

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
			actual := ExactlyOneSet(tt.s, tt.ss...)
			if tt.expected != actual {
				t.Errorf("expected %t; got %t", tt.expected, actual)
			}
		})
	}
}

func TestAge(t *testing.T) {
	tests := []struct {
		name     string
		t        time.Time
		now      time.Time
		expected string
	}{
		{
			name:     "exactly now",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC),
			expected: "just now",
		},
		{
			name:     "within a few milliseconds",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			expected: "just now",
		},
		{
			name:     "10 seconds",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 22, 22, 21, time.UTC),
			expected: "10s",
		},
		{
			name:     "10 minutes",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 32, 12, 21, time.UTC),
			expected: "10m",
		},
		{
			name:     "24 hours",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 18, 15, 22, 12, 21, time.UTC),
			expected: "1d",
		},
		{
			name:     "25 hours",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 18, 16, 22, 12, 21, time.UTC),
			expected: "1d",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := Age(tt.t, tt.now)
			if tt.expected != actual {
				t.Errorf("expected %s; got %s", tt.expected, actual)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	wrapped := Wrap("json", map[string]interface{}{
		"foo": "bar",
	})
	jsonString, _ := json.Marshal(wrapped)
	assert.JSONEq(t, `{"json": {"foo": "bar"}}`, string(jsonString))
}
