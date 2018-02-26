package main

import (
	"net/http"
	"time"
	"strings"
)

func serveWelcome(w http.ResponseWriter, r *http.Request) {
	welcome := string(MustAsset("res/welcome.html"))

	poems := ParsePoems("./poems.txt")
	now := time.Now()
	poemsIndex := now.Day() % len(poems)
	poem := poems[poemsIndex]
	linesIndex := int(now.Weekday()) % len(poem.LinesCode)

	welcome = strings.Replace(welcome, "<PoemTitle/>", poem.Title, 1)
	welcome = strings.Replace(welcome, "<PoemAuthor/>", poem.Author, 1)

	lines := ""
	for i, line := range poem.Lines {
		if i == linesIndex {
			lines += `<div style="color:red">` + line + `</div>`
		} else {
			lines += `<div>` + line + `</div>`
		}
	}

	welcome = strings.Replace(welcome, "<PoemLines/>", lines, 1)

	w.Write([]byte(welcome))
}
