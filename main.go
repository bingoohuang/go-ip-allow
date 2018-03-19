package main

import (
	"github.com/bingoohuang/go-utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc(conf.ContextPath+"/", serveWelcome)
	http.HandleFunc(conf.ContextPath+"/home", go_utils.RandomPoemBasicAuth(serveToolsIndex))
	http.HandleFunc(conf.ContextPath+"/favicon.png", serveFavicon)
	http.HandleFunc(conf.ContextPath+"/ipAllow", go_utils.RandomPoemBasicAuth(serveIpAllow)) // 设置IP权限

	sport := strconv.Itoa(conf.ListenPort)
	log.Println("start to listen at ", sport)

	http.ListenAndServe(":"+sport, nil)
}

func serveWelcome(w http.ResponseWriter, r *http.Request) {
	welcome := MustAsset("res/welcome.html")

	go_utils.ServeWelcome(w, welcome, conf.ContextPath)
}

func serveToolsIndex(w http.ResponseWriter, r *http.Request) {
	forceReset := r.FormValue("forceReset")

	if forceReset == "1" {
		serveHome(w, r, "")
		return
	}

	loginUserName, cookie := login(r)
	log.Println("loginUserName:", loginUserName, ",cookie", cookie)
	msg := ""
	if loginUserName != "" {
		newIpFileLine, err := ipAllow(cookie, loginUserName)
		if err == nil {
			go_utils.ClearCookie(w, conf.CookieName)
			serveTools(cookie.OfficeIp, newIpFileLine, w)
			return
		}
		msg = err.Error()
	}

	alreadySet, clientIp, ipFileLine := isIpAlreadySet(r)
	if !alreadySet {
		serveHome(w, r, msg)
		return
	}

	serveTools(clientIp, *ipFileLine, w)
}

func serveTools(clientIp string, ipFileLine IpFileLine, w http.ResponseWriter) {
	toolsHtml := string(MustAsset("res/tools.html"))
	toolsHtml = strings.Replace(toolsHtml, "{{CLIENT_IP}}", clientIp, -1)
	toolsHtml = strings.Replace(toolsHtml, "{{USER_NAME}}", ipFileLine.User, -1)
	toolsHtml = strings.Replace(toolsHtml, "{{SETT_IME}}", ipFileLine.Day, -1)
	w.Write([]byte(toolsHtml))
}

func isIpAlreadySet(r *http.Request) (bool, string, *IpFileLine) {
	clientIP := go_utils.GetClientIp(r)
	isPrivateIP, _ := go_utils.IsPrivateIP(clientIP)
	if isPrivateIP {
		return false, clientIP, nil
	}

	fileLines := ParseAllowIpsFile()
	for _, fileLine := range fileLines {
		if clientIP == fileLine.Ip {
			return true, clientIP, &fileLine
		}
	}

	return false, clientIP, nil
}
