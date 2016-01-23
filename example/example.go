package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/vitess/go/ioutil2"
	"github.com/gophergala2016/dagger"
)

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
	return ioutil2.WriteFileAtomic(task.Output().Path, b, 0644)
}

// Output to a file.
func (task GithubUser) Output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: fmt.Sprintf("./GithubUser-%s.json", task.Username)}
}

func main() {
	task := GithubUser{Username: "gophergala2016"}
	if !task.Output().Exists() {
		log.Printf("running task: %+v", task)
		if err := task.Run(); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("target done: %s", task.Output().Path)
}
