package handler

import (
	"sort"
	"strings"

	"github.com/notwithering/shit/tree"
)

func groupByDir(children []*tree.Node, less func(i, j int) bool) func(i, j int) bool {
	return func(i, j int) bool {
		if children[i].IsDirectory != children[j].IsDirectory {
			return children[i].IsDirectory
		}
		return less(i, j)
	}
}

func sortName(children []*tree.Node) {
	sort.Slice(children, groupByDir(children, func(i, j int) bool {
		li, lj := strings.ToLower(children[i].Name), strings.ToLower(children[j].Name)
		if li != lj {
			return li < lj
		}
		return children[i].Name < children[j].Name
	}))
}
