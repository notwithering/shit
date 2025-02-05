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

	goFileServerFlag = kingpin.Flag("go", "Use Go's http.FileServer.").Short('g').Bool()
	goFileServer     bool

	uploadFlag = kingpin.Flag("upload", "Allow file uploading.").Short('u').Bool()
	upload     bool

	indexFlag = kingpin.Flag("index", "Automatically serve index files.").Short('i').Bool()
	index     bool

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
	kingpin.Parse()

	host = *hostFlag
	port = *portFlag
	goFileServer = *goFileServerFlag
	upload = *uploadFlag
	index = *indexFlag
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

func checkForFlagIncompatabilities() {
	var err bool

	if useTLS && (tlsCert == "" || tlsKey == "") {
		kingpin.Errorf("flags --cert and --key are required when --tls is set")
		err = true
	}
	if goFileServer {
		if rootMode != rootModeSingleDir {
			kingpin.Errorf("flag --go only compatible with rootModeSingleDir")
			err = true
		}
		if upload {
			kingpin.Errorf("flag --upload incompatible with --go")
			err = true
		}
		if index {
			kingpin.Errorf("flag --index incompatible with --go")
			err = true
		}
	}
	if upload && rootMode == rootModeSingleFile {
		kingpin.Errorf("flag --upload incompatible with rootModeSingleFile")
		err = true
	}

	if err {
		os.Exit(1)
	}
}
