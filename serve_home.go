package main

import (
	"fmt"
	"net/http"
	"strings"
	"log"
)

func serverHome(w http.ResponseWriter, r *http.Request) {
	logined, cookie := login(r)
	log.Println("logined:", logined, ",cookie", cookie)
	msg := ""
	if logined {
		msg = ipAllow(cookie)
	}
	clearCookie(w)

	envCheckboxes := ""
	for _, env := range conf.Envs {
		envCheckboxes += fmt.Sprintf("<input class='env' type='checkbox' checked value='%v'>%v</input><br/>", env, env)
	}

	html := string(MustAsset("res/index.html"))
	js := string(MustAsset("res/index.js"))
	if logined {
		js = strings.Replace(js, "/*.ALERTS*/", `alert('`+msg+`')`, 1)
	}

	html = strings.Replace(html, "<envCheckboxes/>", envCheckboxes, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)

	w.Write([]byte(html))
}
