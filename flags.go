package main

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

var (
	hostFlag = kingpin.Flag("host", "The host to bind to.").Short('h').Default("0.0.0.0").Envar("HOST").String()
	host     string

	portFlag = kingpin.Flag("port", "The port to serve.").Short('p').Default("8080").Envar("PORT").String()
	port     string

	uploadFlag = kingpin.Flag("upload", "Allow file uploading.").Short('u').Bool()
	upload     bool

	useTLSFlag = kingpin.Flag("tls", "Enable TLS.").Short('t').Bool()
	useTLS     bool

	tlsCertFlag = kingpin.Flag("cert", "Path to TLS certificate file.").Short('c').Envar("TLS_CERT").ExistingFile()
	tlsCert     string

	tlsKeyFlag = kingpin.Flag("key", "Path to the TLS key file.").Short('k').Envar("TLS_KEY").ExistingFile()
	tlsKey     string

	exportsArg = kingpin.Arg("files", "The files or directories to share.").Default(".").ExistingFilesOrDirs()
	exports    []string
)

func parseFlags() {
	var err error
	tmpl, err = parseTemplate()
	if err != nil {
		kingpin.Fatalf("error while parsing html template: %s", err)
	}

	kingpin.Parse()

	host = *hostFlag
	port = *portFlag
	upload = *uploadFlag
	useTLS = *useTLSFlag
	tlsCert = *tlsCertFlag
	tlsKey = *tlsKeyFlag

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
}
