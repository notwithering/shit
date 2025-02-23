package main

var (
	indexFiles = []string{"index.html", "default.html", "index.htm", "home.html", "default.htm", "index.php", "default.php"}
)

const (
	rootModeSingleDir int = iota
	rootModeExports
	rootModeSingleFile
)
