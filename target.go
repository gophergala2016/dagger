package dagger

import "os"

type BooleanTarget struct {
	Value bool
}

func (t BooleanTarget) Exists() bool {
	return t.Value
}

// LocalTarget represents a file on the file system.
type LocalTarget struct {
	Path string
}

// Exists returns, whether this path exists.
func (t LocalTarget) Exists() bool {
	if _, err := os.Stat(t.Path); os.IsNotExist(err) {
		return false
	}
	return true
}
