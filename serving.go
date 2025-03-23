package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	rootMode int
)

func startServer() {
	s := &http.Server{
		Addr:         host + ":" + port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	if useTLS {
		fmt.Printf("https://%s:%s/\n", host, port)
		if err := s.ListenAndServeTLS(tlsCert, tlsKey); err != nil {
			kingpin.Fatalf("error while serving: %s", err)
		}
	} else {
		fmt.Printf("http://%s:%s/\n", host, port)
		if err := s.ListenAndServe(); err != nil {
			kingpin.Fatalf("error while serving: %s", err)
		}
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) error {
	export := exports[0]

	if index {
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

	if rootMode == rootModeSingleDir {
		dir, err := os.ReadDir(export)
		if err != nil {
			return err
		}

		return executeDirEntries(w, r, dir)
	} else if rootMode == rootModeSingleFile {
		http.Redirect(w, r, filepath.Base(export), http.StatusPermanentRedirect)
		return nil
	}

	return executeExports(w, r)
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
			http.Redirect(w, r, r.URL.Path+"/", http.StatusPermanentRedirect)
			return nil
		}

		if index {
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

		return executeDirEntries(w, r, dir)
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
