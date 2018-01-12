package hcloud

// String returns a pointer to the passed string s.
func String(s string) *string { return &s }

// Int returns a pointer to the passed integer i.
func Int(i int) *int { return &i }
