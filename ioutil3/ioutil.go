package ioutil3

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"code.google.com/p/vitess/go/ioutil2"
)

func WriteFileAtomic(filename string, data []byte, perm os.FileMode) error {
	dirname := filepath.Dir(filename)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return err
		}
	}
	return ioutil2.WriteFileAtomic(filename, data, perm)
}

// WriteTabs like io.WriteString, just for a list of values.
func WriteTabs(w io.Writer, record ...string) error {
	return WriteTabsDelimiter(w, '\t', record...)
}

// WriteTabs like io.WriteString, just for a list of values.
func WriteTabsDelimiter(w io.Writer, delim rune, record ...string) error {
	wr := csv.NewWriter(w)
	wr.Comma = delim
	if err := wr.Write(record); err != nil {
		return err
	}
	wr.Flush()
	return wr.Error()
}

func WriteJSON(w io.Writer, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if _, err = w.Write(b); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}
