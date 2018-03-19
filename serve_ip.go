package main

import (
	"github.com/bingoohuang/go-utils"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func serveIpAllow(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	officeIp := strings.TrimSpace(req.FormValue("officeIp"))
	if !go_utils.IsIP4Valid(officeIp) {
		w.Write([]byte(`IP格式非法`))
		return
	}

	envs := strings.TrimSpace(req.FormValue("envs"))
	log.Println("officeIp:", officeIp, ",env:", envs)

	csrfToken := go_utils.RandString(10)
	go_utils.WriteCookie(w, conf.EncryptKey, conf.CookieName, &CookieValue{
		OfficeIp:  officeIp,
		Envs:      envs,
		CsrfToken: csrfToken,
		Expired:   time.Now().Add(time.Duration(24) * time.Hour),
	})
	url := go_utils.CreateWxQyLoginUrl(conf.WxCorpId, conf.WxAgentId, conf.RedirectUri, csrfToken)
	log.Println("wx login url:", url)

	w.Write([]byte(url))
}

func ipAllow(cookie *CookieValue, loginUserName string) (IpFileLine, error) {
	allowedIpLines := CreateAllowIpsFileLines(cookie.Envs, cookie.OfficeIp, loginUserName)
	log.Println("CreateAllowIpsFileLines:", allowedIpLines)

	allowedIps := joinAllowedIpLines(allowedIpLines)

	go func() {
		out, err := exec.Command("/bin/bash", conf.UpdateFirewallShell, cookie.Envs, allowedIps).Output()
		if err != nil {
			log.Println("设置失败，执行SHELL错误", err)
			//return allowedIpLines[0], errors.New(`设置失败，执行SHELL错误` + err.Error())
		}

		shellOut := string(out)
		log.Println(shellOut)
	}()

	saveAllowIpsFile(allowedIpLines)
	return allowedIpLines[0], nil
}
