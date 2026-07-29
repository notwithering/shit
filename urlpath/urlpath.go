package urlpath

import (
	"strings"
)

func Clean(path string) string {
	path = strings.TrimLeft(path, "/")
	isDir := strings.HasSuffix(path, "/")
	path = strings.TrimRight(path, "/")

	slice := make([]string, 0, len(path))

	for part := range strings.SplitSeq(path, "/") {
		switch part {
		case "", ".":
			continue
		case "..":
			if len(slice) > 0 {
				slice = slice[:len(slice)-1]
			}
		default:
			slice = append(slice, part)
		}
	}

	if len(slice) == 0 {
		return "/"
	}

	result := "/" + strings.Join(slice, "/")
	if isDir {
		result += "/"
	}

	return result
}

func Join(path ...string) string {
	return Clean(strings.Join(path, "/"))
}

func SplitList(path string) []string {
	path = strings.TrimLeft(path, "/")
	path = strings.TrimRight(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

func SetDir(path string, isDir bool) string {
	if isDir {
		return Clean(path + "/")
	}
	return Clean(strings.TrimRight(path, "/"))
}
