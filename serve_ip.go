package main

import (
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func serveIpAllow(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	officeIp := strings.TrimSpace(req.FormValue("officeIp"))
	if !isIP4Valid(officeIp) {
		w.Write([]byte(`IP格式非法`))
		return
	}

	envs := strings.TrimSpace(req.FormValue("envs"))
	log.Println("officeIp:", officeIp, ",env:", envs)

	csrfToken := RandomString(10)

	writeCsrfTokenCookie(w, csrfToken, officeIp, envs)
	url := createWxQyLoginUrl(conf.RedirectUri, csrfToken)
	log.Println("wx login url:", url)

	w.Write([]byte(url))
}

func ipAllow(cookie *CookieValue) string {
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

func isIP4Valid(ipv4 string) bool {
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(ipv4)
}
