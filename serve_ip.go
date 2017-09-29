package main

import (
	"io/ioutil"
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
	url := createWxQyLoginUrl(g_config.RedirectUri, csrfToken)
	log.Println("wx login url:", url)

	w.Write([]byte(url))
}

func ipAllow(r *http.Request, cookie *CookieValue) string {
	_, allowedIps := isIpAlreadyAllowed(cookie.Envs, cookie.OfficeIp)

	if allowedIps == "" {
		allowedIps = cookie.OfficeIp
	} else {
		allowedIps += "," + cookie.OfficeIp
	}

	out, err := exec.Command("/bin/bash", g_config.UpdateFirewallShell, cookie.Envs, allowedIps).Output()
	if err != nil {
		return `设置失败，执行SHELL错误` + err.Error()
	}

	shellOut := string(out)
	log.Println(shellOut)
	writeAllowIpFile(cookie.Envs, allowedIps)

	return `设置成功`
}

func isIpAlreadyAllowed(envs, ip string) (bool, string) {
	content, err := ioutil.ReadFile(envs + "-AllowIps.txt")
	if err != nil {
		return false, ""
	}

	strContent := string(content)
	alreadyAllowed := strings.Contains(strContent, ip)
	return alreadyAllowed, strContent
}

func writeAllowIpFile(env, content string) {
	ioutil.WriteFile(env+"-AllowIps.txt", []byte(content), 0644)
}
