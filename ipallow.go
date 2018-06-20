package main

import (
	"github.com/bingoohuang/go-utils"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func serveIpAllow(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	officeIp := strings.TrimSpace(r.FormValue("officeIp"))
	if !go_utils.IsIP4Valid(officeIp) {
		w.Write([]byte(`IP格式非法`))
		return
	}

	envs := strings.TrimSpace(r.FormValue("envs"))
	log.Println("officeIp:", officeIp, ",env:", envs)

	cookie := r.Context().Value("CookieValue").(*go_utils.CookieValueImpl)
	ipAllow(cookie, envs, officeIp)

	w.Write([]byte("OK"))
}

func ipAllow(cookie *go_utils.CookieValueImpl, envs, officeIp string) {
	allowedIpLines := CreateAllowIpsFileLines(envs, officeIp, cookie.Name)
	allowedIps := joinAllowedIpLines(allowedIpLines)

	go func() {
		out, err := exec.Command("/bin/bash", conf.UpdateFirewallShell, envs, allowedIps).Output()
		if err != nil {
			log.Println("设置失败，执行SHELL错误", err)
		}

		shellOut := string(out)
		log.Println(shellOut)
	}()

	saveAllowIpsFile(allowedIpLines)
}
