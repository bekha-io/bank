package utils

func Contains(s string, ss []string) bool {
	for _, el := range ss {
		if el == s {
			return true
		}
	}
	return false
}
