package dagger

import "log"

type input struct {
	deps []interface{}
}

// Input returns a dispatch object for a given task, that
// provides many helper methods for dealing with
func Input(in interface{}) input {
	switch t := in.(type) {
	case Requirer:
		reqs := t.Requires()
		log.Printf("task implements Requirer: %+v, %T", reqs, reqs)
		switch tt := reqs.(type) {
		case Outputter:
			log.Printf("task %v has a single dependency: %v", tt, reqs)
		case []Outputter:
			log.Println("list of outputters ...")
			for i, v := range tt {
				log.Printf("%d. %v, %T", i, v, v)
			}
		case []interface{}:
			log.Println("list of interface ...")
			for i, v := range tt {
				log.Printf("%d. %v, %T", i, v, v)
			}
		default:
			log.Printf("something else: %s", tt)
		}
	default:
		log.Printf("all done: task %v has no requirements", in)
	}
	return input{}
}
