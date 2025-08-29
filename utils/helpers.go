package utils

func PickNonEmptyString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
