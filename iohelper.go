package dagger

// TaskMap expressed the dependencies of a task. It should be easy to create
// and should provide many helper methods.
type TaskMap map[string]Outputter

// Path return the path to the output for a given key. If no key is found or
// the outputter is not a LocalTarget, the empty string is returned.
func (tm TaskMap) Path(key string) string {
	task, ok := tm[key]
	if !ok {
		return ""
	}
	output := task.Output()
	switch o := output.(type) {
	case LocalTarget:
		return o.Path
	default:
		return ""
	}
}
