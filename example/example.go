package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"code.google.com/p/vitess/go/ioutil2"
	"github.com/gophergala2016/dagger"
)

type Executable struct {
	Name string
}

func (task Executable) Output() dagger.Target {
	_, err := exec.LookPath(task.Name)
	return dagger.BooleanTarget{Value: err == nil}
}

// GithubUser downloads the user info.
type GithubUser struct {
	Username string
}

// Run downloads the information.
func (task GithubUser) Run() error {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s", task.Username))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil2.WriteFileAtomic(task.Output().(dagger.LocalTarget).Path, b, 0644)
}

// Output to a file.
func (task GithubUser) Output() dagger.Target {
	return dagger.LocalTarget{Path: fmt.Sprintf("./GithubUser-%s.json", task.Username)}
}

// GithubRepos gets the list of repositories for a given user.
type GithubRepos struct {
	Username string
}

func (task GithubRepos) Requires() interface{} {
	return []interface{}{
		GithubUser{Username: task.Username},
		Executable{Name: "jq"},
	}
}

func (task GithubRepos) Run() error {
	return nil
}

// Output to a file.
func (task GithubRepos) Output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: fmt.Sprintf("./GithubRepos-%s.json", task.Username)}
}

// makeDeps return the predecessors of a task.
func makeDeps(task interface{}) (deps []dagger.Outputter) {
	switch t := task.(type) {
	case dagger.Requirer:
		log.Printf("%+v, %T", t, t)
	default:
		log.Println("task has no requirements")
	}
	return
}

func main() {
	task := GithubRepos{Username: "gophergala2016"}
	dagger.Input(task)
	// log.Println(makeDeps(task))
	// jq := Executable{Name: "jq"}
	// log.Println(jq.Output().Exists())

	// task := GithubRepos{Username: "gophergala2016"}
	// if !task.Output().Exists() {
	// 	log.Printf("running task: %+v", task)
	// 	if err := task.Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// log.Printf("target done: %s", task.Output().Path)
}
