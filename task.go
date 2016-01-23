package dagger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/fatih/structs"
)

// Target is the result of a task. It can exist or not.
type Target interface {
	Exists() bool
}

// Outputter is something, that can result in a target.
type Outputter interface {
	Output() Target
}

// Runner is something to be executed. Business logic goes here.
type Runner interface {
	Run() error
}

// Requirer is implemented by tasks which have requirements. The dependencies
// are passed as map as it is the most flexible: Single dependency, dependency
// list or a map of dependencies.
type Requirer interface {
	Requires() TaskMap
}

// TaskMap expressed the dependencies of a task. It should be easy to create
// and should provide many helper methods.
type TaskMap map[string]Outputter

// inputDispatcher provides shortcuts to let a task access its requirements.
type inputDispatcher struct {
	r Requirer
}

func Input(r Requirer) inputDispatcher {
	return inputDispatcher{r: r}
}

func (d inputDispatcher) Scanner() (*bufio.Scanner, error) {
	for _, v := range d.r.Requires() {
		output := v.Output()
		switch o := output.(type) {
		case LocalTarget:
			file, err := o.Open()
			if err != nil {
				return nil, err
			}
			return bufio.NewScanner(file), nil
		}
	}
	return nil, nil
}

func (d inputDispatcher) ReadLines() (<-chan string, error) {
	ch := make(chan string)
	for _, v := range d.r.Requires() {
		output := v.Output()
		switch o := output.(type) {
		case LocalTarget:
			file, err := o.Open()
			if err != nil {
				return nil, err
			}
			go func() {
				reader := bufio.NewReader(file)
				for {
					line, err := reader.ReadString('\n')
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatal(err)
					}
					ch <- strings.TrimSpace(line)
				}
				close(ch)
			}()
		default:
			close(ch)
		}
	}
	return ch, nil
}

// AutoPath returns a path based on the task name and parameters.
func AutoPath(outp Outputter) string {
	return AutoPathExt(outp, "tsv")
}

// AutoPath returns a path based on the task name and parameters.
func AutoPathExt(outp Outputter, ext string) string {
	directory := strings.Replace(fmt.Sprintf("%T", outp), ".", "/", -1)
	m := structs.Map(outp)

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, k)
		value := m[k]
		switch v := value.(type) {
		case string:
			parts = append(parts, v)
		case fmt.Stringer:
			parts = append(parts, v.String())
		default:
			parts = append(parts, fmt.Sprintf("%v", m[k]))
		}
	}
	fn := strings.Join(parts, "-")
	if len(fn) == 0 {
		fn = "output"
	}
	filename := fn + "." + ext
	return path.Join(directory, filename)
}
