package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gophergala2016/dagger"
	"github.com/kr/pretty"
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
	file, err := task.output().Create()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.WriteString(file, fmt.Sprintf("%d\n", sum))
	return err
}

// output points to a file.
func (task BogusAggregation) output() dagger.LocalTarget {
	return dagger.LocalTarget{Path: dagger.AutoPathExt(task, "tsv")}
}

// Output fulfills the Outputter interface.
func (task BogusAggregation) Output() dagger.Target {
	return task.output()
}

type Dummy struct{}

func (task Dummy) Requires() dagger.TaskMap {
	return dagger.TaskMap{
		"x": BogusData{Length: 133, Cols: 10},
		"y": BogusAggregation{Col: 999},
	}
}

func (task Dummy) Run() error {
	return nil
}

func (task Dummy) Output() dagger.Target {
	return dagger.LocalTarget{Path: dagger.AutoPath(task)}
}

func main() {
	// task := BogusAggregation{Col: 2}
	task := Dummy{}

	prereqs := dagger.TaskDeps(task)
	log.Printf("%# v", pretty.Formatter(prereqs))
	for i, o := range dagger.TopoSort(prereqs) {
		log.Printf("%d: %# v - %# v - %v", i, o, o.Output(), o.Output().Exists())
		if !o.Output().Exists() {
			log.Printf("running %# v...", o)
			if rr, ok := o.(dagger.Runner); ok {
				if err := dagger.InitializeRunner(&rr); err != nil {
					log.Fatal(err)
				}
				if err := rr.Run(); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal("cannot create missing task output for: %# v", o)
			}
		}
	}
	// for _, v := range task.Requires() {
	// 	output := v.Output()
	// 	if !output.Exists() {
	// 		switch tt := v.(type) {
	// 		case dagger.Runner:
	// 			if err := tt.Run(); err != nil {
	// 				log.Fatal(err)
	// 			}
	// 		default:
	// 			log.Fatal("output does not exists and there is nothing to run")
	// 		}
	// 	} else {
	// 		log.Printf("dependency is done")
	// 	}
	// }

	// if !task.Output().Exists() {
	// 	if err := task.Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// log.Println("all done")
}
