package main

import (
	"fmt"

	"github.com/notwithering/shit/tree"
)

func main() {
	kc, cli := parseCli()

	tree, err := tree.GenerateTree(cli.Exports)
	kc.FatalIfErrorf(err, "generating tree")
	if tree == nil {
		kc.Fatalf("no valid exports")
	}

	fmt.Println(tree)

	err = listenAndServeTree(cli, tree)
	kc.FatalIfErrorf(err, "starting server")
}
