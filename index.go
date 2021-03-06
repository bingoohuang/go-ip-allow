package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strings"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	envCheckboxes := ""
	for _, env := range conf.Envs {
		envCheckboxes += fmt.Sprintf("<input class='env' type='checkbox' checked value='%v'>%v</input><br/>", env, env)
	}

	html := string(MustAsset("res/index.html"))
	html = strings.Replace(html, "<envCheckboxes/>", envCheckboxes, 1)
	html = go_utils.MinifyHtml(html, false)
	html = strings.Replace(html, "${ContextPath}", conf.ContextPath, -1)
	clientIP := go_utils.GetClientIp(r)
	isPrivateIP, _ := go_utils.IsPrivateIP(clientIP)
	if isPrivateIP {
		clientIP = "识别中请稍待,或请拷贝下面的IP后设置。"
	}

	html = strings.Replace(html, "/*.AUTOIP*/", clientIP, 1)

	w.Write([]byte(html))
}
