package strings

import "unicode"

// Capitalize capitalizes the first letter of a string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// IsCapitalized checks if a string starts with an uppercase letter.
func IsCapitalized(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}
