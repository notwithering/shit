package main

import (
	"fmt"
	"io/fs"
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

	if rootMode == rootModeSingleDir {
		dir, err := os.ReadDir(export)
		if err != nil {
			return err
		}

		return executeDirEntries(w, dir)
	} else if rootMode == rootModeSingleFile {
		http.Redirect(w, r, filepath.Base(export), http.StatusPermanentRedirect)
		return nil
	}

	return executeExports(w)
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

		dir, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		return executeDirEntries(w, dir)
	}

	return serveFile(w, r, path)
}

func executeExports(w http.ResponseWriter) error {
	var links []string

	for _, file := range exports {
		fileinfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		name := filepath.Base(file)
		if fileinfo.IsDir() {
			name += "/"
		}
		links = append(links, filepath.Base(file))
	}

	return tmpl.Execute(w, tmplData{
		Links:  links,
		Upload: upload,
	})
}

func executeDirEntries(w http.ResponseWriter, dir []fs.DirEntry) error {
	var links []string

	for _, file := range dir {
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		links = append(links, name)
	}

	return tmpl.Execute(w, tmplData{
		Links:  links,
		Upload: upload,
	})
}
