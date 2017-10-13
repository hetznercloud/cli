package cli

func yesno(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
