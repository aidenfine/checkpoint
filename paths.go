package checkpoint

import "strings"

func MatchPathPattern(pattern, path string) bool {
	if pattern == path {
		return true
	}

	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")

		// does path start with the prefix?
		if !strings.HasPrefix(path, parts[0]) {
			return false
		}

		if len(parts) == 2 {
			return strings.HasSuffix(path, parts[1])
		}
		currentIndex := len(parts[0])
		for i := 1; i < len(parts)-1; i++ {
			if parts[i] == "" {
				continue
			}

			index := strings.Index(path[currentIndex:], parts[i])
			if index == -1 {
				return false
			}

			currentIndex += index + len(parts[i])
		}
		lastPart := parts[len(parts)-1]
		if lastPart != "" {
			return strings.HasSuffix(path, lastPart)
		}
		return true
	}
	return false
}
