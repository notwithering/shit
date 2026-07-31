package tree

import (
	"errors"
)

func GenerateTree(exports []string) (*Node, error) {
	if len(exports) == 0 {
		return nil, errors.New("no exports")
	}

	if len(exports) == 1 {
		return NodeFromPath(exports[0])
	}

	root := &Node{
		Name:        "",
		IsDirectory: true,
	}

	for _, path := range exports {
		child, err := NodeFromPath(path)
		if err != nil {
			return nil, err
		}

		root.Children = append(root.Children, child)
	}

	return root, nil
}
