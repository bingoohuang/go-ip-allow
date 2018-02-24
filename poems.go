package main

import (
	"io/ioutil"
	"fmt"
	"strings"
)

type Poem struct {
	Title     string
	TitleCode string
	Author    string
	Lines     []string
	LinesCode []string
}

func ParsePoems(poemFile string) []Poem {
	poems := make([]Poem, 0)
	poemsBytes, err := ioutil.ReadFile(poemFile)
	if err != nil {
		fmt.Println("read poems error", err.Error())
		return poems
	}

	fileLines := strings.Split(string(poemsBytes), "\n")

	for i := 0; i < len(fileLines); i++ {
		l := strings.TrimSpace(fileLines[i])

		if l == "" {
			continue
		}

		titleFields := strings.SplitN(l, "#", 2)
		i++
		author := strings.TrimSpace(fileLines[i])

		lines := make([]string, 0)
		linesCode := make([]string, 0)
		for i++; i < len(fileLines); i++ {
			if fileLines[i] == "" {
				break
			}

			lineFields := strings.SplitN(fileLines[i], "#", 2)
			lines = append(lines, lineFields[0])
			linesCode = append(linesCode, lineFields[1])
		}

		poems = append(poems, Poem{
			Title:     titleFields[0],
			TitleCode: titleFields[1],
			Author:    author,
			Lines:     lines,
			LinesCode: linesCode,
		})
	}

	return poems
}
