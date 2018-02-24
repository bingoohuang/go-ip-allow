package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc(conf.ContextPath+"/", serveWelcome)
	http.HandleFunc(conf.ContextPath+"/home", BasicAuth(serveHome))
	http.HandleFunc(conf.ContextPath+"/favicon.png", serveFavicon)
	http.HandleFunc(conf.ContextPath+"/ipAllow", BasicAuth(serveIpAllow)) // 设置IP权限

	sport := strconv.Itoa(conf.ListenPort)
	log.Println("start to listen at ", sport)

	http.ListenAndServe(":"+sport, nil)
}
