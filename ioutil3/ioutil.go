package ioutil3

import (
	"encoding/csv"
	"encoding/json"
	"io"
)

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

func WriteJSON(w io.Writer, o interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
