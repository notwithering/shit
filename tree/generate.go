package tree

import (
	"errors"
	"os"
	"path/filepath"
)

func GenerateTree(exports []string) (*Node, error) {
	if len(exports) == 0 {
		return nil, errors.New("no exports")
	}

	if len(exports) == 1 {
		path, err := filepath.Abs(exports[0])
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		return &Node{
			Name:        info.Name(),
			IsDirectory: info.IsDir(),
			Physical:    true,
			Path:        path,
		}, nil
	}

	root := &Node{
		Name:        "",
		IsDirectory: true,
		Physical:    false,
	}

	for _, path := range exports {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		child := &Node{
			Name:        info.Name(),
			IsDirectory: info.IsDir(),
			Physical:    true,
			Path:        path,
		}

		root.Children = append(root.Children, child)
	}

	return root, nil
}
