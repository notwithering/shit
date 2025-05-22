package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func registerHandlers() {
	if cli.Go {
		http.Handle("/", http.FileServer(http.Dir(cli.Exports[0])))
	} else {
		registerShitHandlers()
	}
}

func registerShitHandlers() {
	http.HandleFunc("GET /", get)
	if cli.Upload {
		http.HandleFunc("POST /", post)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if err := serveRoot(w, r); err != nil {
			httpErr(w, err)
		}

		return
	}

	if err := serveExports(w, r); err != nil {
		httpErr(w, err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	rc.SetReadDeadline(time.Now().Add(cli.UploadTimeout))

	if err := r.ParseMultipartForm(int64(cli.MaxUploadMemory)); err != nil {
		httpErr(w, err)
		return
	}

	path, err := getRealPath(r.URL.Path)
	if err != nil {
		httpErr(w, err)
		return
	}

	if path == "" {
		http.NotFound(w, r)
		return
	}

	fileinfo, err := os.Stat(path)
	if err != nil {
		httpErr(w, err)
		return
	}

	if !fileinfo.IsDir() {
		http.NotFound(w, r)
		return
	}

	files := r.MultipartForm.File["files"]
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			httpErr(w, err)
			return
		}
		defer file.Close()

		dst, err := os.Create(filepath.Join(path, header.Filename))
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			httpErr(w, err)
			return
		}

		if rootMode == rootModeExports {
			cli.Exports = append(cli.Exports, header.Filename)
		}
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
