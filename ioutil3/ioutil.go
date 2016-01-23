package ioutil3

import (
	"encoding/csv"
	"io"
)

// WriteTabs like io.WriteString, just for a list of values.
func WriteTabs(w io.Writer, record []string) error {
	return WriteTabsDelimiter(w, record, '\t')
}

// WriteTabs like io.WriteString, just for a list of values.
func WriteTabsDelimiter(w io.Writer, record []string, delim rune) error {
	wr := csv.NewWriter(w)
	wr.Comma = delim
	if err := wr.Write(record); err != nil {
		return err
	}
	wr.Flush()
	return wr.Error()
}
