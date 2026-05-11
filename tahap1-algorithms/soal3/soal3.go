package soal3

// ValidateBrackets checks if a string of bracket characters is valid.
// Uses a manual stack implementation; no regex allowed.
func ValidateBrackets(s string) bool {
	// Length check
	if len(s) < 1 || len(s) > 4096 {
		return false
	}

	// Valid characters
	valid := map[byte]bool{
		'<': true, '>': true,
		'{': true, '}': true,
		'[': true, ']': true,
	}

	// Matching pairs: closer -> expected opener
	pair := map[byte]byte{
		'>': '<',
		'}': '{',
		']': '[',
	}

	var stack []byte

	for i := 0; i < len(s); i++ {
		ch := s[i]

		// Check valid character
		if !valid[ch] {
			return false
		}

		if pair[ch] == 0 {
			// Opening bracket
			stack = append(stack, ch)
		} else {
			// Closing bracket
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if top != pair[ch] {
				return false
			}
		}
	}

	return len(stack) == 0
}
