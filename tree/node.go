package tree

import "strings"

type Node struct {
	Name        string
	IsDirectory bool
	Children    []*Node

	Physical bool
	Path     string
}

func (n *Node) String() string {
	var sb strings.Builder

	sb.WriteString(n.Name)
	if n.IsDirectory {
		sb.WriteString("/")
	}
	if !n.Physical {
		sb.WriteString(" (virtual)")
	}
	sb.WriteString("\n")

	for i, child := range n.Children {
		isLast := i == len(n.Children)-1
		connector := "├─"
		continuation := "│ "
		if isLast {
			connector = "└─"
			continuation = "  "
		}
		var j int
		for line := range strings.SplitSeq(child.String(), "\n") {
			if j == 0 {
				sb.WriteString(connector)
			} else {
				sb.WriteString(continuation)
			}

			sb.WriteString(line)
			sb.WriteString("\n")

			j++
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
