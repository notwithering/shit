package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/notwithering/shit/tree"
	"github.com/notwithering/shit/urlpath"
)

func FileServer(root *tree.Node) http.Handler {
	return &fileHandler{root: root}
}

type fileHandler struct {
	root *tree.Node
}

func (h *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := urlpath.Clean(r.URL.Path)

	node, err := h.pathToNode(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if node == nil {
		http.Error(w, "404: Not Found", http.StatusNotFound)
		return
	}

	path = urlpath.SetDir(path, node.IsDirectory)

	if path != r.URL.Path {
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	serveNode(w, r, node)
}

func (h *fileHandler) pathToNode(path string) (*tree.Node, error) {
	current := h.root
	pathParts := urlpath.SplitList(path)

pathWalker:
	for partI, part := range pathParts {
		if current.Physical {
			rest := pathParts[partI:]
			if len(rest) == 0 {
				// fully resolved
				return current, nil
			}

			if !current.IsDirectory {
				// cant resolve further than a file
				return nil, nil
			}
			realPath := filepath.Join(append([]string{current.Path}, rest...)...) // FIXME

			info, err := os.Stat(realPath)
			if err != nil {
				if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ENOTDIR) {
					return nil, nil
				}
				return nil, err
			}

			return &tree.Node{
				Name:        info.Name(),
				IsDirectory: info.IsDir(),
				Physical:    true,
				Path:        realPath,
			}, nil
		}

		for _, child := range current.Children {
			if child.Name == part {
				current = child
				continue pathWalker
			}
		}

		return nil, nil
	}

	return current, nil
}

func serveNode(w http.ResponseWriter, r *http.Request, node *tree.Node) {
	if node.IsDirectory {
		serveDirectory(w, r, node)
	} else {
		serveFile(w, r, node)
	}
}

func serveDirectory(w http.ResponseWriter, _ *http.Request, n *tree.Node) {
	var children []*tree.Node

	if n.Physical {
		dirEntries, err := os.ReadDir(n.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, e := range dirEntries {
			children = append(children, &tree.Node{
				Name:        e.Name(),
				IsDirectory: e.IsDir(),
				Physical:    true,
				Path:        filepath.Join(n.Path, e.Name()),
			})
		}
	} else {
		children = n.Children
	}

	for _, child := range children {
		name := child.Name
		if child.IsDirectory {
			name += "/"
		}

		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br>", name, name)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, node *tree.Node) {
	http.ServeFile(w, r, node.Path)
}
