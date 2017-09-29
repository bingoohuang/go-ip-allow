package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc(g_config.ContextPath+"/", serverHome)
	http.HandleFunc(g_config.ContextPath+"/favicon.png", serveFavicon)
	http.HandleFunc(g_config.ContextPath+"/ipAllow", serveIpAllow) // 设置IP权限

	sport := strconv.Itoa(g_config.ListenPort)
	log.Println("start to listen at ", sport)

	http.ListenAndServe(":"+sport, nil)
}
