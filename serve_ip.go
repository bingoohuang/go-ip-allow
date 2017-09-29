package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func serveIpAllow(w http.ResponseWriter, req *http.Request) {
	officeIp := strings.TrimSpace(req.FormValue("officeIp"))
	envs := strings.TrimSpace(req.FormValue("envs"))
	log.Println("officeIp:", officeIp, ",env:", envs)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	csrfToken := RandomString(10)

	writeCsrfTokenCookie(w, csrfToken, officeIp, envs)
	url := createWxQyLoginUrl(conf.RedirectUri, csrfToken)
	log.Println("wx login url:", url)

	w.Write([]byte(url))
}

func ipAllow(r *http.Request, cookie *CookieValue) string {
	allowedIpLines := parseAllowIpsFile(cookie.Envs, cookie.OfficeIp)
	allowedIps := joinAllowedIpLines(allowedIpLines)

	out, err := exec.Command("/bin/bash", conf.UpdateFirewallShell, cookie.Envs, allowedIps).Output()
	if err != nil {
		return `设置失败，执行SHELL错误` + err.Error()
	}

	shellOut := string(out)
	log.Println(shellOut)
	saveAllowIpsFile(cookie.Envs, allowedIpLines)

	return `设置成功`
}
