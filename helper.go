package main

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func httpErr(w http.ResponseWriter, err error) {
	if !errors.Is(err, http.ErrHandlerTimeout) {
		kctx.Errorf("[%s] %s\n", time.Now().Format(time.DateTime), err.Error())
	}
	if w.Header().Get("Content-Type") == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRealPath(path string) (string, error) {
	path = strings.TrimLeft(path, "/")
	path = filepath.Clean(path)

	getFile := func(export string, file string) (string, error) {
		path := filepath.Join(export, file)

		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		} else if err != nil {
			return "", err
		}

		relPath, err := filepath.Rel(export, path)
		if err != nil {
			return "", err
		}

		if !strings.HasPrefix(relPath, "..") {
			return path, nil
		}

		return "", nil
	}

	if rootMode == rootModeSingleDir {
		for _, export := range cli.Exports {
			file, err := getFile(export, path)
			if err != nil {
				return "", nil
			}
			if file != "" {
				return file, nil
			}
		}
	} else {
		split := strings.Split(path, string(filepath.Separator))
		reqExport := split[0]

		if reqExport == "." { // r.URL.Path == "/"
			return ".", nil
		}

		reqFile := filepath.Join(split[1:]...)

		for _, export := range cli.Exports {
			if filepath.Base(export) == reqExport {
				file, err := getFile(export, reqFile)
				if err != nil {
					kctx.Errorf("error while finding requested file: %s", err)
				}

				return file, nil
			}
		}
	}

	return "", nil
}

func redirectCode() int {
	if cli.PermanentRedirect {
		return http.StatusPermanentRedirect
	}
	return http.StatusTemporaryRedirect
}
