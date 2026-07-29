package main

import (
	"time"

	"github.com/alecthomas/kong"
)

type cli struct {
	Host string `help:"The host to bind to." short:"h" env:"HOST" default:"0.0.0.0"`
	Port string `help:"The port to serve." short:"p" env:"PORT" default:"8080"`

	TLS  bool   `help:"Enable TLS." short:"t"`
	Cert string `help:"Path to TLS certificate file." short:"c" env:"TLS_CERT" type:"existingfile"`
	Key  string `help:"Path to TLS key file." short:"k" env:"TLS_KEY" type:"existingfile"`

	ReadTimeout  time.Duration `help:"Timeout for a request to complete." default:"5s"`
	WriteTimeout time.Duration `help:"Timeout for a response to complete." default:"10s"`

	Exports []string `arg:"" name:"files" help:"The files or directories to share." type:"existingfileexistingdir" default:"."`
}

func parseCli() (*kong.Context, *cli) {
	var cli cli
	kc := kong.Parse(&cli)
	return kc, &cli
}
