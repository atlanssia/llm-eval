package embed

import (
	"io/fs"
)

// FS returns the embedded filesystem for the React frontend
// Note: This is a placeholder - actual embedding is done in root embed.go
func FS() (fs.FS, error) {
	return nil, fs.ErrNotExist
}

// Exists returns true if the embedded frontend exists
func Exists() bool {
	return false
}
