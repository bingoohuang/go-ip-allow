package main

import (
	"fmt"
	"net/http"
	"strings"
)

func serverHome(w http.ResponseWriter, r *http.Request) {
	ok, cookie := login(r)
	msg := ""
	if ok {
		msg = ipAllow(r, cookie)
	}
	clearCookie(w)

	envCheckboxes := ""
	for _, env := range conf.Envs {
		envCheckboxes += fmt.Sprintf("<input class='env' type='checkbox' checked value='%v'>%v</input><br/>", env, env)
	}

	html := string(MustAsset("res/index.html"))
	js := string(MustAsset("res/index.js"))
	if ok {
		js = strings.Replace(js, "/*.ALERTS*/", `alert('`+msg+`')`, 1)
	}

	html = strings.Replace(html, "<envCheckboxes/>", envCheckboxes, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)

	w.Write([]byte(html))
}
