package utils

// DerefString returns a string from *string
// or empty string if pointer is nil
func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}
