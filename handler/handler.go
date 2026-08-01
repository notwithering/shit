package handler

import (
	"fmt"
	"html"
	"net/http"
	"path/filepath"

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
		internalServerError(w, err)
		return
	}
	if node == nil {
		notFound(w, r)
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

			node, err := tree.NodeFromPath(realPath)
			if err != nil {
				return nil, err
			}

			return node, nil
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
	w.Header().Set("Cache-Control", "no-cache")

	if node.IsDirectory {
		serveDirectory(w, r, node)
	} else {
		serveFile(w, r, node)
	}
}

func serveDirectory(w http.ResponseWriter, r *http.Request, n *tree.Node) {
	var children []*tree.Node

	if n.Physical {
		var err error
		children, err = n.ReadDir()
		if err != nil {
			internalServerError(w, err)
			return
		}
		if children == nil {
			notFound(w, r)
			return
		}
	} else {
		children = n.Children
	}

	sortName(children)

	w.WriteHeader(http.StatusOK)

	for _, child := range children {
		name := child.Name
		if child.IsDirectory {
			name += "/"
		}
		name = html.EscapeString(name)

		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br>", name, name)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, node *tree.Node) {
	if !node.Physical {
		http.Error(w, "virtual file somehow", http.StatusInternalServerError)
		return
	}

	info, err := node.Stat()
	if err != nil {
		internalServerError(w, err)
		return
	}
	if info == nil {
		notFound(w, r)
		return
	}

	file, err := node.Open()
	if err != nil {
		internalServerError(w, err)
		return
	}
	if file == nil {
		notFound(w, r)
		return
	}
	defer file.Close()

	http.ServeContent(w, r, node.Name, info.ModTime(), file)
}
