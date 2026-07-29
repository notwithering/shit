package main

import (
	"net/http"

	"github.com/notwithering/shit/handler"
	"github.com/notwithering/shit/tree"
)

func listenAndServeTree(cli *cli, root *tree.Node) error {
	s := &http.Server{
		Addr:         cli.Host + ":" + cli.Port,
		ReadTimeout:  cli.ReadTimeout,
		WriteTimeout: cli.WriteTimeout,
	}

	s.Handler = handler.FileServer(root)

	if cli.TLS && cli.Cert != "" && cli.Key != "" {
		return s.ListenAndServeTLS(cli.Cert, cli.Key)
	}
	return s.ListenAndServe()
}
