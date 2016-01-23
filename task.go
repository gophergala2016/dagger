package dagger

import (
	"fmt"
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
		parts = append(parts, fmt.Sprintf("%s", m[k]))
	}
	filename := strings.Join(parts, "-") + "." + ext
	return path.Join(directory, filename)
}
