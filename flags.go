package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
)

var cli struct {
	Host              string        `help:"The host to bind to." short:"h" env:"HOST" default:"0.0.0.0"`
	Port              string        `help:"The port to serve." short:"p" env:"PORT" default:"8080"`
	Go                bool          `help:"Use Go's http.FileServer." short:"g"`
	Upload            bool          `help:"Allow file uploading." short:"u"`
	MaxUploadMemory   ByteSize      `help:"The maximum memory allowed when saving uploaded files." short:"m" default:"10MiB"`
	Index             bool          `help:"Automatically serve index files." short:"i"`
	TLS               bool          `help:"Enable TLS." short:"t"`
	Cert              string        `help:"Path to TLS certificate file." short:"c" env:"TLS_CERT" type:"existingfile"`
	Key               string        `help:"Path to TLS key file." short:"k" env:"TLS_KEY" type:"existingfile"`
	ReadTimeout       time.Duration `help:"Timeout for a request to complete." short:"r" default:"5s"`
	WriteTimeout      time.Duration `help:"Timeout for a response to complete." short:"w" default:"10s"`
	UploadTimeout     time.Duration `help:"Timeout for a file upload to complete." short:"U" default:"30m"`
	PermanentRedirect bool          `help:"Use permanent redirects." short:"P"`
	Exports           []string      `arg:"" name:"files" help:"The files or directories to share." type:"existingfileexistingdir" default:"."`
}

var kctx *kong.Context

func parseFlags() {
	kctx = kong.Parse(&cli)

	var exports []string
	for _, export := range cli.Exports {
		abs, err := filepath.Abs(export)
		kctx.FatalIfErrorf(err, "error while getting export %s's absolute path", export)

		for _, e := range exports {
			if filepath.Base(e) == filepath.Base(abs) {
				kctx.Fatalf("can't have 2 exports with same base name: %s, %s", e, abs)
			}
		}
		exports = append(exports, abs)
	}
	cli.Exports = exports

	if len(exports) == 1 {
		info, err := os.Stat(exports[0])
		kctx.FatalIfErrorf(err, "error while finding root mode")

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
	var hasErr bool

	if cli.TLS && (cli.Cert == "" || cli.Key == "") {
		println("flags --cert and --key are required when --tls is set")
		hasErr = true
	}
	if cli.Go {
		if rootMode != rootModeSingleDir {
			kctx.Errorf("flag --go only compatible with rootModeSingleDir")
			hasErr = true
		}
		if cli.Upload {
			kctx.Errorf("flag --upload incompatible with --go")
			hasErr = true
		}
		if cli.Index {
			kctx.Errorf("flag --index incompatible with --go")
			hasErr = true
		}
		if cli.PermanentRedirect {
			kctx.Errorf("flag --permanent-redirect incompatible with --go")
			hasErr = true
		}
	}
	if cli.Upload && rootMode == rootModeSingleFile {
		kctx.Errorf("flag --upload incompatible with rootModeSingleFile")
		hasErr = true
	}

	if hasErr {
		os.Exit(1)
	}
}

type ByteSize int64

func (b *ByteSize) Decode(ctx *kong.DecodeContext) error {
	var raw string
	if err := ctx.Scan.PopValueInto("size", &raw); err != nil {
		return err
	}
	bytes, err := humanize.ParseBytes(raw)
	if err != nil {
		return err
	}
	*b = ByteSize(bytes)
	return nil
}
