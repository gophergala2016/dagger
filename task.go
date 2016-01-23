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

// Inputter is implemented by tasks which have requirements. The dependencies
// are passed as map as it is the most flexible: Single dependency, dependency
// list or a map of dependencies.
type Inputter interface {
	Input() TaskMap
}
