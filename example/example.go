package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"code.google.com/p/vitess/go/ioutil2"
	"github.com/gophergala2016/dagger"
	"github.com/gophergala2016/dagger/ioutil3"
)

// Executable is a thin task, making sure required programs are installed.
type Executable struct {
	Name string
}

// Output will exist, if the executable is in PATH.
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
	return ioutil2.WriteFileAtomic(task.output().Path, b, 0644)
}

func (task GithubUser) output() dagger.LocalTarget {
	// TODO(miku): add support for automatic nameing, like: dagger.Autolocated(task)
	return dagger.LocalTarget{Path: fmt.Sprintf("GithubUser-%s.json", task.Username)}
}

// Output to a file, conformance to interface.
func (task GithubUser) Output() dagger.Target {
	return task.output()
}

// GithubRepos gets the list of repositories for a given user.
type GithubRepos struct {
	Username string
}

func (task GithubRepos) Requires() dagger.TaskMap {
	return dagger.TaskMap{
		"user": GithubUser{Username: task.Username},
		"jq":   Executable{Name: "jq"},
	}
}

func (task GithubRepos) Run() error {
	file, err := task.output().Create()
	if err != nil {
		return err
	}
	defer file.Close()
	if err := ioutil3.WriteTabs(file, []string{"Hello", "World", "1", "2"}); err != nil {
		return err
	}
	return nil
}

func (task GithubRepos) output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: fmt.Sprintf("./GithubRepos-%s.json", task.Username)}
}

// Output to a file.
func (task GithubRepos) Output() dagger.Target {
	return task.output()
}

func main() {
	task := GithubRepos{Username: "gophergala2016"}
	log.Printf("%+v", task.Requires())
	if err := task.Run(); err != nil {
		log.Fatal(err)
	}
}
