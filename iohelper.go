package dagger

// TaskMap expressed the dependencies of a task. It should be easy to create
// and should provide many helper methods.
type TaskMap map[string]Outputter
