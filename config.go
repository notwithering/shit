package main

import (
	"time"
)

var (
	indexFiles = []string{"index.html", "index.htm", "index.php", "index.md", "default.html"}
)

const (
	readTimeout           = 5 * time.Second
	writeTimeout          = 10 * time.Second
	maxUploadMemory int64 = 10 << 20 // 10 MiB

	rootModeSingleDir int = iota
	rootModeExports
	rootModeSingleFile
)
