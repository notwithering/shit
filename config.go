package main

import (
	"time"
)

var (
	indexFiles = []string{"index.html", "default.html", "index.htm", "home.html", "default.htm", "index.php", "default.php"}
)

const (
	readTimeout           = 5 * time.Second
	writeTimeout          = 10 * time.Second
	maxUploadMemory int64 = 10 << 20 // 10 MiB

	rootModeSingleDir int = iota
	rootModeExports
	rootModeSingleFile
)
