package soal1

import "strings"

// FindMatchingStrings finds case-insensitive duplicate strings and returns
// their 1-based indices of the first matching set found.
// Returns nil, false if no duplicates exist.
func FindMatchingStrings(n int, strs []string) ([]int, bool) {
	if n <= 1 || len(strs) != n {
		return nil, false
	}

	// Track every occurrence of each lowercase string (1-based indices)
	occurrences := make(map[string][]int)
	for i := 0; i < n; i++ {
		key := strings.ToLower(strs[i])
		occurrences[key] = append(occurrences[key], i+1)
	}

	// Pick the group with the highest frequency.
	// If multiple groups tie, the one whose first occurrence appears
	// first in the input wins.
	var result []int
	maxLen := 0
	bestStart := n + 1

	for _, indices := range occurrences {
		if len(indices) >= 2 {
			if len(indices) > maxLen || (len(indices) == maxLen && indices[0] < bestStart) {
				result = indices
				maxLen = len(indices)
				bestStart = indices[0]
			}
		}
	}

	if len(result) == 0 {
		return nil, false
	}
	return result, true
}