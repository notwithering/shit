package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
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

	rootModeSingleDir int = iota
	rootModeExports
	rootModeSingleFile
)

var (
	tmpl = template.Must(template.New("").Parse(tmplStr))

	portFlag   = kingpin.Flag("port", "The port to serve.").Short('p').Default("8080").String()
	port       string
	exportsArg = kingpin.Arg("files", "The files or directories to share.").Default(".").ExistingFilesOrDirs()
	exports    []string

	rootMode int
)

func main() {
	kingpin.Parse()
	port = *portFlag
	for _, export := range *exportsArg {
		abs, err := filepath.Abs(export)
		if err != nil {
			kingpin.Fatalf("error while getting export %s's absolute path: %s", export, err)
		}
		for _, export := range exports {
			if filepath.Base(export) == filepath.Base(abs) {
				kingpin.Fatalf("can't have 2 exports with same base name: %s, %s", export, abs)
			}
		}
		exports = append(exports, abs)
	}

	if len(exports) == 1 {
		info, err := os.Stat(exports[0])
		if err != nil {
			kingpin.Fatalf("error while finding root mode: %s", err)
		}

		if info.IsDir() {
			rootMode = rootModeSingleDir
		} else {
			rootMode = rootModeSingleFile
		}
	} else {
		rootMode = rootModeExports
	}

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
		kingpin.Fatalf("error while serving: %s", err)
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
	path := filepath.Clean(r.URL.Path[1:])
	var file string

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
		for _, export := range exports {
			var err error
			file, err = getFile(export, path)
			if err != nil {
				kingpin.Errorf("error while finding requested file: %s", err)
				continue
			}
			if file != "" {
				break
			}
		}
	} else {
		split := strings.Split(path, string(filepath.Separator))
		reqExport := split[0]
		reqFile := filepath.Join(split[1:]...)

		for _, export := range exports {
			if filepath.Base(export) == reqExport {
				var err error
				file, err = getFile(export, reqFile)
				if err != nil {
					kingpin.Errorf("error while finding requested file: %s", err)
				}

				break
			}
		}
	}

	if file == "" {
		http.NotFound(w, r)
		return nil
	}

	fileinfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	if fileinfo.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusPermanentRedirect)
			return nil
		}

		dir, err := os.ReadDir(file)
		if err != nil {
			return err
		}

		return executeDirEntries(w, dir)
	}

	return serveFile(w, file)
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
	kingpin.Errorf("[%s] %s\n", time.Now().Format(time.DateTime), err.Error())
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
