package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func serverHome(w http.ResponseWriter, r *http.Request) {
	logined, cookie := login(r)
	log.Println("logined:", logined, ",cookie", cookie)
	msg := ""
	if logined {
		msg = ipAllow(cookie)
	}
	clearCookie(w)

	envCheckboxes := ""
	for _, env := range conf.Envs {
		envCheckboxes += fmt.Sprintf("<input class='env' type='checkbox' checked value='%v'>%v</input><br/>", env, env)
	}

	html := string(MustAsset("res/index.html"))
	js := string(MustAsset("res/index.js"))
	if logined {
		js = strings.Replace(js, "/*.ALERTS*/", `alert('`+msg+`')`, 1)
	}

	html = strings.Replace(html, "<envCheckboxes/>", envCheckboxes, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)

	w.Write([]byte(html))
}

func login(r *http.Request) (bool, *CookieValue) {
	loginCookie := readLoginCookie(r)
	if loginCookie == nil {
		return false, nil
	}

	ok, _ := tryLogin(loginCookie, r)
	return ok, loginCookie
}

func tryLogin(loginCookie *CookieValue, r *http.Request) (bool, error) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	log.Println("code:", code, ",state:", state)
	if loginCookie != nil && code != "" && state == loginCookie.CsrfToken {
		accessToken, err := getAccessToken(conf.WxCorpId, conf.WxCorpSecret)
		if err != nil {
			return false, err
		}
		userId, err := getLoginUserId(accessToken, code)
		if err != nil {
			return false, err
		}
		userInfo, err := getUserInfo(accessToken, userId)
		if err != nil {
			return false, err
		}

		sendLoginInfo(userInfo, loginCookie, accessToken)

		return true, nil
	}

	return false, errors.New("no login")
}

func sendLoginInfo(info *WxUserInfo, loginCookie *CookieValue, accessToken string) {
	msg := map[string]interface{}{
		"touser":  "@all",
		"toparty": "@all",
		"totag":   "@all",
		"msgtype": "text",
		"agentid": conf.WxAgentId,
		"text": map[string]string{
			"content": "用户[" + info.Name + "]正在扫码设置IP[" + loginCookie.OfficeIp + "]，环境[" + loginCookie.Envs + "]",
		},
		"safe": 0,
	}
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken
	_, err := httpPost(url, msg)
	if err != nil {
		log.Println("sendLoginInfo error", err)
	}
}
