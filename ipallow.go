package main

import (
	"github.com/bingoohuang/go-utils"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"encoding/json"
)

func serveIpAllow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	officeIp := strings.TrimSpace(r.FormValue("officeIp"))
	if !go_utils.IsIP4Valid(officeIp) {
		w.Write([]byte(`IP格式非法`))
		return
	}

	envs := strings.TrimSpace(r.FormValue("envs"))
	log.Println("officeIp:", officeIp, ",env:", envs)

	cookie := CookieValue{}
	cookieValue := r.Header.Get("CookieValue")
	json.Unmarshal([]byte(cookieValue), &cookie)

	ipAllow(&cookie, envs, officeIp)

	w.Write([]byte("OK"))
}

func ipAllow(cookie *CookieValue, envs, officeIp string) {
	allowedIpLines := CreateAllowIpsFileLines(envs, officeIp, cookie.Name)
	log.Println("CreateAllowIpsFileLines:", allowedIpLines)

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
