package tree

import (
	"errors"
)

func GenerateTree(exports []string) (*Node, error) {
	var nodes []*Node

	for _, export := range exports {
		node, err := NodeFromPath(export)
		if err != nil {
			return nil, err
		}
		if node == nil {
			continue
		}

		nodes = append(nodes, node)
	}

	if len(nodes) == 0 {
		return nil, errors.New("no exports")
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	return &Node{
		Name:        "",
		IsDirectory: true,
		Children:    nodes,
	}, nil
}
