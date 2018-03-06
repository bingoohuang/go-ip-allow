package main

import (
	"net/http"
	"strings"
	"github.com/bingoohuang/go-utils"
)

func serveWelcome(w http.ResponseWriter, r *http.Request) {
	welcome := string(MustAsset("res/welcome.html"))

	poem, linesIndex := go_utils.RandomPoem()

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
