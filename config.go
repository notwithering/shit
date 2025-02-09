package main

import (
	"time"
)

var (
	indexFiles = []string{"index.html", "default.html", "index.htm", "home.html", "default.htm", "index.php", "default.php"}
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second

	rootModeSingleDir int = iota
	rootModeExports
	rootModeSingleFile
)
