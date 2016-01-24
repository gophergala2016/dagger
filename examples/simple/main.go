package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gophergala2016/dagger"
)

// BogusData creates some bogus data.
type BogusData struct {
	Length int `default:"100"`
	Cols   int `default:"3"`
}

// Run creates the data.
func (task BogusData) Run() error {
	file, err := dagger.Output(task).CreateLocalTarget()
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < task.Length; i++ {
		var vals []string
		for j := 0; j < task.Cols; j++ {
			vals = append(vals, strconv.Itoa(rand.Intn(1000)))
		}
		if _, err := io.WriteString(file, strings.Join(vals, "\t")+"\n"); err != nil {
			return err
		}
	}

	return nil
}

// Output fulfills the Outputter interface.
func (task BogusData) Output() dagger.Target {
	return dagger.LocalTarget{Path: dagger.AutoPathExt(task, "tsv")}
}

// BogusAggregation
type BogusAggregation struct {
	Col int `default:"1"`
}

func (task BogusAggregation) Requires() dagger.TaskMap {
	return dagger.TaskMap{
		"data":  BogusData{Length: 1000, Cols: task.Col + 1},
		"train": BogusData{Length: 100, Cols: task.Col + 1},
	}
}

func (task BogusAggregation) Run() error {
	scanner, err := dagger.Input(task).Scanner()
	if err != nil {
		return err
	}
	var sum int
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		v, err := strconv.Atoi(fields[task.Col])
		if err != nil {
			return err
		}
		sum += v
	}
	file, err := dagger.Output(task).CreateLocalTarget()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.WriteString(file, fmt.Sprintf("%d\n", sum))
	return err
}

// Output fulfills the Outputter interface.
func (task BogusAggregation) Output() dagger.Target {
	return dagger.LocalTarget{Path: dagger.AutoPathExt(task, "tsv")}
}

type Dummy struct {
	Length int `default:"133"`
}

// Requires returns some things this task requires.
func (task Dummy) Requires() dagger.TaskMap {
	return dagger.TaskMap{
		"x": BogusData{Length: task.Length, Cols: 10},
		"y": BogusAggregation{Col: 999},
	}
}

// Run tries to be short.
func (task Dummy) Run() error {
	file := dagger.Output(task).MustCreateLocalTarget()
	defer file.Close()
	for _, p := range dagger.Input(task).PathList() {
		if _, err := io.WriteString(file, fmt.Sprintf("OK\t%s\n", p)); err != nil {
			return err
		}
	}
	return nil
}

// Output.
func (task Dummy) Output() dagger.Target {
	return dagger.LocalTarget{Path: dagger.AutoPath(task)}
}

func main() {
	task := Dummy{}
	if err := dagger.Build(task); err != nil {
		log.Fatal(err)
	}
}
