package util

import "golang.org/x/net/idna"

// DisplayZoneName prepares the zone name for user display by converting the ASCII punycode encoding to Unicode.
func DisplayZoneName(zoneName string) string {
	display, err := idna.Display.ToUnicode(zoneName)
	if err != nil {
		// Let's just show the zoneName as returned by the API, better than throwing an error.
		return zoneName
	}
	return display
}

// ParseZoneIDOrName converts user input for the zone name into an API compatible form. We want to accept Unicode and
// Punycode/ASCII variants from the user, but the API always expects Punycode/ASCII.
func ParseZoneIDOrName(zoneIDOrName string) (string, error) {
	return idna.Punycode.ToASCII(zoneIDOrName)
}
