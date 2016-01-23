package dagger

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

// Dependency is implemented by tasks, which require only a single other task.
type Dependency interface {
	Requires() Outputter
}

// DependencyList is implemented by tasks, which require a list of things to be done.
type DependencyList interface {
	Requires() []Outputter
}

// DependencyMap is implemented by tasks, which require a map of tasks to be done.
type DependencyMap interface {
	Requires() map[string]Outputter
}

// Task is simply something that produces an output and knows, how to
// fabricate it, if if does not exist.
type Task interface {
	Runner
	Outputter
}
