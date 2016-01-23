package dagger

import (
	"os"
	"path/filepath"

	"github.com/facebookgo/atomicfile"
)

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

func (t LocalTarget) Create() (*atomicfile.File, error) {
	dirname := filepath.Dir(t.Path)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return nil, err
		}
	}
	return atomicfile.New(t.Path, 0644)
}

func (t LocalTarget) Open() (*os.File, error) {
	return os.Open(t.Path)
}
