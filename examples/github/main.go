package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/cenkalti/backoff"
	"github.com/google/go-github/github"
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
	return ioutil3.WriteFileAtomic(task.output().Path, b, 0644)
}

func (task GithubUser) output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: dagger.AutoPathExt(task, "json")}
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

// Run downloads all repos for a given username and save as line delimited JSON.
func (task GithubRepos) Run() error {
	client := github.NewClient(nil)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}
	var allRepos []github.Repository

	for {
		var (
			repos []github.Repository
			resp  *github.Response
			err   error
		)
		operation := func() error {
			repos, resp, err = client.Repositories.ListByOrg(task.Username, opt)
			log.Println(err)
			if err != nil {
				return err
			}
			return nil
		}
		if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	file, err := task.output().Create()
	if err != nil {
		return err
	}
	defer file.Close()
	for _, repo := range allRepos {
		if err := ioutil3.WriteJSON(file, repo); err != nil {
			return err
		}
	}
	return nil
}

// output for internal user.
func (task GithubRepos) output() dagger.LocalTarget {
	return dagger.LocalTarget{dagger.AutoPath(task)}
}

// Output to a file.
func (task GithubRepos) Output() dagger.Target {
	return task.output()
}

func main() {
	task := GithubUser{Username: "gophergala2016"}
	// task := GithubRepos{Username: "gophergala2016"}
	if !task.Output().Exists() {
		if err := task.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("all done")
	}
}
