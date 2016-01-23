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

func (task GithubRepos) Input() dagger.TaskMap {
	return dagger.TaskMap{
		"user": GithubUser{Username: task.Username},
		"jq":   Executable{Name: "jq"},
	}
}

func (task GithubRepos) Run() error {
	log.Println(task.Input().Path("user"))
	return nil
}

// Output to a file.
func (task GithubRepos) Output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: fmt.Sprintf("./GithubRepos-%s.json", task.Username)}
}

func main() {
	task := GithubRepos{Username: "gophergala2016"}
	log.Printf("%+v", task.Input())
	if err := task.Run(); err != nil {
		log.Fatal(err)
	}
}
