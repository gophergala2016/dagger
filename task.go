package dagger

// Target is the result of a task. It can exist or not.
type Target interface {
	Exists() bool
}

// Outputter is something, that can result in a target.
type Outputter interface {
	Output() Target
}

// Completer is another way for a task to specify, whether it is complete or not.
type Completer interface {
	Complete() bool
}

// Runner is something to be executed. Business logic goes here.
type Runner interface {
	Run() error
}

// Dependency is implemented by tasks, which require one or more things to be
// done, before they can run.
type Requirer interface {
	Requires() interface{}
}

// Task is simply something that produces an output and knows, how to
// fabricate it, if if does not exist.
type Task interface {
	Runner
	Outputter
}
