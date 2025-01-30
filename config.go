package main

import (
	"time"
)

const (
	readTimeout           = 5 * time.Second
	writeTimeout          = 10 * time.Second
	maxUploadMemory int64 = 10 << 20 // 10MB

	rootModeSingleDir int = iota
	rootModeExports
	rootModeSingleFile
)
