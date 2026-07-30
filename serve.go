package main

import (
	"net/http"

	"github.com/notwithering/shit/handler"
	"github.com/notwithering/shit/tree"
)

func listenAndServeTree(cli *cli, root *tree.Node) error {
	s := &http.Server{
		Addr:    cli.Host + ":" + cli.Port,
		Handler: handler.FileServer(root),

		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 0,
	}

	if cli.TLS && cli.Cert != "" && cli.Key != "" {
		return s.ListenAndServeTLS(cli.Cert, cli.Key)
	}
	return s.ListenAndServe()
}
