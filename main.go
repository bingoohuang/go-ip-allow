package main

import (
	"log"
	"net/http"
	"strconv"
	"github.com/bingoohuang/go-utils"
)

func main() {
	http.HandleFunc(conf.ContextPath+"/", serveWelcome)
	http.HandleFunc(conf.ContextPath+"/home", go_utils.RandomPoemBasicAuth(serveHome))
	http.HandleFunc(conf.ContextPath+"/favicon.png", serveFavicon)
	http.HandleFunc(conf.ContextPath+"/ipAllow", go_utils.RandomPoemBasicAuth(serveIpAllow)) // 设置IP权限

	sport := strconv.Itoa(conf.ListenPort)
	log.Println("start to listen at ", sport)

	http.ListenAndServe(":"+sport, nil)
}
