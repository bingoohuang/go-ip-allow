package main

import (
	"io/ioutil"
	"fmt"
	"strings"
)

type Poem struct {
	Title string
	TitleCode string
	Author string
	Lines []string
	LinesCode []string
}

func ParsePoems() []Poem {
	poems := make([]Poem, 0)
	poemsBytes, err := ioutil.ReadFile("./poems.txt")
	if err != nil {
		fmt.Println("read poems error", err.Error())
		return poems
	}

	lines := strings.Split(string(poemsBytes), "\n")
	var poem *Poem = nil

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			if poem != nil {
				poems = append(poems, *poem)
				poem = nil
			}
		} else {
			strings.SplitN(l, "#", 2)
			
			if poem == nil {
				poem = &Poem {
					Title:l,
				}
			}
		}
	}

	return poems
}
