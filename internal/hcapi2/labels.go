package hcapi2

func labelKeys(m map[string]string) []string {
	ks := make([]string, len(m))
	i := 0
	for k := range m {
		ks[i] = k
		i++
	}
	return ks
}
