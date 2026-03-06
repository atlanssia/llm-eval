package main

import (
	"embed"
	"io/fs"
)

//go:embed web/dist
var webFS embed.FS

// GetWebFS returns the embedded filesystem for the React frontend
func GetWebFS() (fs.FS, error) {
	return fs.Sub(webFS, "web/dist")
}

// WebFSExists returns true if the embedded frontend exists
func WebFSExists() bool {
	_, err := webFS.ReadDir("web/dist")
	return err == nil
}
