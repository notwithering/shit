package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	rootMode int
)

func startServer() {
	s := &http.Server{
		Addr:         cli.Host + ":" + cli.Port,
		ReadTimeout:  cli.ReadTimeout,
		WriteTimeout: cli.WriteTimeout,
	}

	if cli.UseTLS {
		fmt.Printf("https://%s:%s/\n", cli.Host, cli.Port)
		if err := s.ListenAndServeTLS(cli.TLSCert, cli.TLSKey); err != nil {
			kctx.Fatalf("error while serving: %s", err)
		}
	} else {
		fmt.Printf("http://%s:%s/\n", cli.Host, cli.Port)
		if err := s.ListenAndServe(); err != nil {
			kctx.Fatalf("error while serving: %s", err)
		}
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) error {
	export := cli.Exports[0]

	if cli.Index {
		path, err := getRealPath(r.URL.Path)
		if err != nil {
			return err
		}

		served, err := serveIndex(w, r, path)
		if err != nil {
			return err
		}
		if served {
			return nil
		}
	}

	var links []string

	if rootMode == rootModeSingleDir {
		dir, err := os.ReadDir(export)
		if err != nil {
			return err
		}

		links = dirToLinks(dir)
	} else if rootMode == rootModeSingleFile {
		http.Redirect(w, r, filepath.Base(export), redirectCode())
		return nil
	} else {
		var err error
		links, err = exportsToLinks(cli.Exports)
		if err != nil {
			return err
		}
	}

	return execute(w, r, links)
}

func serveExports(w http.ResponseWriter, r *http.Request) error {
	path, err := getRealPath(r.URL.Path)
	if err != nil {
		return err
	}

	if path == "" {
		http.NotFound(w, r)
		return nil
	}

	fileinfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileinfo.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", redirectCode())
			return nil
		}

		if cli.Index {
			served, err := serveIndex(w, r, path)
			if err != nil {
				return err
			}
			if served {
				return nil
			}
		}

		dir, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		links := dirToLinks(dir)
		execute(w, r, links)
	}

	return serveFile(w, r, path)
}

func serveIndex(w http.ResponseWriter, r *http.Request, path string) (bool, error) {
	for _, indexFile := range indexFiles {
		indexPath := filepath.Join(path, indexFile)
		if _, err := os.Stat(indexPath); err == nil {
			return true, serveFile(w, r, indexPath)
		}
	}
	return false, nil
}

func serveFile(w http.ResponseWriter, r *http.Request, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	fileinfo, err := os.Stat(filename)
	if err != nil {
		return err
	}

	http.ServeContent(w, r, filename, fileinfo.ModTime(), file)
	return nil
}
