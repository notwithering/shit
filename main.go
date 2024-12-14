package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/h2non/filetype"
)

const (
	defaultPort = "8080"
)

var (
	tmpl = template.Must(template.New("").Parse(tmplStr))

	portFlag   = kingpin.Flag("port", "The port to serve.").Short('p').Default("8080").String()
	port       string
	exportsArg = kingpin.Arg("files", "The files or directories to share.").Default(".").ExistingFilesOrDirs()
	exports    []string
)

func main() {
	kingpin.Parse()
	port = *portFlag
	exports = *exportsArg

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if err := serveRoot(w, r); err != nil {
				httpErr(w, err)
			}

			return
		}

		if err := serveExports(w, r); err != nil {
			httpErr(w, err)
		}
	})

	fmt.Printf("http://127.0.0.1:%s/\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) error {
	if len(exports) == 1 {
		export := exports[0]

		fileinfo, err := os.Stat(export)
		if err != nil {
			return err
		}

		if fileinfo.IsDir() {
			dir, err := os.ReadDir(export)
			if err != nil {
				return err
			}

			return executeDirEntries(w, dir)
		} else {
			http.Redirect(w, r, filepath.Base(export), http.StatusPermanentRedirect)
			return nil
		}
	}

	return executeExports(w)
}

func serveExports(w http.ResponseWriter, r *http.Request) error {
	requestedFile := filepath.Clean(r.URL.Path[1:])
	var filePath string

	for _, export := range exports {
		fileinfo, err := os.Stat(export)
		if err != nil {
			return err
		}

		if !fileinfo.IsDir() {
			export = "."
		}

		file := filepath.Join(export, requestedFile)

		fileinfo, err = os.Stat(file)
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return err
		}

		exportAbs, err := filepath.Abs(export)
		if err != nil {
			return err
		}

		fileAbs, err := filepath.Abs(export)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(exportAbs, fileAbs)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(relPath, "..") {
			filePath = file
			break
		}
	}

	if filePath == "" {
		http.NotFound(w, r)
		return nil
	}

	fileinfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if fileinfo.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusPermanentRedirect)
			return nil
		}

		dir, err := os.ReadDir(filePath)
		if err != nil {
			return err
		}

		return executeDirEntries(w, dir)
	}

	return serveFile(w, filePath)
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

	return tmpl.Execute(w, links)
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

	return tmpl.Execute(w, links)
}

func httpErr(w http.ResponseWriter, err error) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.DateTime), err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func serveFile(w http.ResponseWriter, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	mtype, err := filetype.Match(buf)
	w.Header().Set("Content-Type", mtype.MIME.Type)

	_, err = w.Write(buf)
	return err
}

const tmplStr string = `<!DOCTYPE html>
<html data-fbscriptallow="true"><head><meta name="viewport" content="width=device-width"></head><body>
<pre>{{range .}}<a href="{{.}}">{{.}}</a>
{{end}}</pre>
</body></html>`
