package ioutil3

import (
	"fmt"
	"io"
	"strings"
)

// WriteTabs like io.WriteString, just for a list of values.
func WriteTabs(w io.Writer, records []interface{}) (int, error) {
	return WriteTabsDelimiter(w, records, "\t")
}

// WriteTabs like io.WriteString, just for a list of values.
func WriteTabsDelimiter(w io.Writer, records []interface{}, delim string) (int, error) {
	var values = make([]string, len(records))
	for i, r := range records {
		values[i] = fmt.Sprintf("%v", r)
	}
	return io.WriteString(w, strings.Join(values, delim))
}
