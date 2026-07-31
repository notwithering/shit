package tree

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type Node struct {
	Name        string
	IsDirectory bool
	Children    []*Node

	Physical bool
	Path     string
}

func (n *Node) Open() (*os.File, error) {
	file, err := os.Open(n.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return file, nil
}

func (n *Node) Stat() (os.FileInfo, error) {
	info, err := os.Stat(n.Path)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ENOTDIR) {
			return nil, nil
		}
		return nil, err
	}

	return info, nil
}

func (n *Node) ReadDir() ([]*Node, error) {
	var children []*Node

	dirEntries, err := os.ReadDir(n.Path)
	if err != nil {
		return children, err
	}

	for _, e := range dirEntries {
		child, err := NodeFromPath(filepath.Join(n.Path, e.Name()))
		if err != nil {
			return children, err
		}

		children = append(children, child)
	}

	return children, nil
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

func NodeFromPath(path string) (*Node, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ENOTDIR) {
			return nil, nil
		}
		return nil, err
	}

	return &Node{
		Name:        info.Name(),
		IsDirectory: info.IsDir(),
		Physical:    true,
		Path:        path,
	}, nil
}
