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
	file, err := task.output().Create()
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

// output points to a file.
func (task BogusData) output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: dagger.AutoPathExt(task, "tsv")}
}

// Output fulfills the Outputter interface.
func (task BogusData) Output() dagger.Target {
	return task.output()
}

// BogusAggregation
type BogusAggregation struct {
	Col int `default:"1"`
}

func (task BogusAggregation) Requires() dagger.TaskMap {
	return dagger.TaskMap{
		"data": BogusData{Length: 1000, Cols: task.Col + 1},
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
	file, err := task.output().Create()
	if err != nil {
		return err
	}
	defer file.Close()
	io.WriteString(file, fmt.Sprintf("%d\n", sum))
	return nil
}

// output points to a file.
func (task BogusAggregation) output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: dagger.AutoPathExt(task, "tsv")}
}

// Output fulfills the Outputter interface.
func (task BogusAggregation) Output() dagger.Target {
	return task.output()
}

func main() {

	task := BogusAggregation{Col: 2}

	for _, v := range task.Requires() {
		output := v.Output()
		if !output.Exists() {
			switch tt := v.(type) {
			case dagger.Runner:
				if err := tt.Run(); err != nil {
					log.Fatal(err)
				}
			default:
				log.Fatal("output does not exists and there is nothing to run")
			}
		} else {
			log.Printf("dependency is done")
		}
	}

	if !task.Output().Exists() {
		if err := task.Run(); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("all done")
}
