package main

import "net/http"

func serveWelcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(MustAsset("res/welcome.html")))
}
