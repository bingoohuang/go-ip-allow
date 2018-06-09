package main

import (
	"github.com/bingoohuang/go-utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
)

func main() {
	http.HandleFunc(conf.ContextPath+"/", MustAuth(serveToolsIndex))
	http.HandleFunc(conf.ContextPath+"/favicon.png", go_utils.ServeFavicon("res/favicon.png", MustAsset, AssetInfo))
	http.HandleFunc(conf.ContextPath+"/ipAllow", MustAuth(serveIpAllow)) // 设置IP权限

	sport := strconv.Itoa(conf.ListenPort)
	log.Println("start to listen at ", sport)

	http.ListenAndServe(":"+sport, nil)
}

func MustAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := CookieValue{}
		go_utils.ReadCookie(r, conf.EncryptKey, conf.CookieName, &cookie)
		if cookie.Name != "" {
			json, _ := json.Marshal(&cookie)
			r.Header.Set("CookieValue", string(json))
			fn(w, r) // 执行被装饰的函数
			return
		}

		http.Redirect(w, r, conf.RedirectUri, 302)
	}
}

func serveToolsIndex(w http.ResponseWriter, r *http.Request) {
	alreadySet, clientIp, ipFileLine := isIpAlreadySet(r)
	if !alreadySet {
		serveHome(w, r)
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
