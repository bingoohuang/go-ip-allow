package main

import (
	"errors"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"log"
	"net/http"
	"strings"
	"time"
)

func serveHome(w http.ResponseWriter, r *http.Request, msg string) {
	envCheckboxes := ""
	for _, env := range conf.Envs {
		envCheckboxes += fmt.Sprintf("<input class='env' type='checkbox' checked value='%v'>%v</input><br/>", env, env)
	}

	js := string(MustAsset("res/index.js"))
	if msg != "" {
		js = strings.Replace(js, "/*.ALERTS*/", `alert('`+msg+`')`, 1)
	}

	js = go_utils.MinifyJs(js, false)

	html := string(MustAsset("res/index.html"))
	html = strings.Replace(html, "<envCheckboxes/>", envCheckboxes, 1)
	html = go_utils.MinifyHtml(html, false)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${ContextPath}", conf.ContextPath, -1)
	clientIP := go_utils.GetClientIp(r)
	isPrivateIP, _ := go_utils.IsPrivateIP(clientIP)
	if isPrivateIP {
		clientIP = "识别中请稍待,或请拷贝下面的IP后设置。"
	}

	html = strings.Replace(html, "/*.AUTOIP*/", clientIP, 1)

	w.Write([]byte(html))
}

type CookieValue struct {
	OfficeIp  string
	Envs      string
	CsrfToken string
	Expired   time.Time
}

func (t *CookieValue) ExpiredTime() time.Time {
	return t.Expired
}

func login(r *http.Request) (string, *CookieValue) {
	cookieValue := &CookieValue{}
	err := go_utils.ReadCookie(r, conf.EncryptKey, conf.CookieName, cookieValue)
	if err != nil {
		return "", nil
	}

	loginUserName, _ := tryLogin(cookieValue, r)
	return loginUserName, cookieValue
}

func tryLogin(loginCookie *CookieValue, r *http.Request) (string, error) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	log.Println("code:", code, ",state:", state)
	if loginCookie != nil && code != "" && state == loginCookie.CsrfToken {
		accessToken, err := go_utils.GetAccessToken(conf.WxCorpId, conf.WxCorpSecret)
		if err != nil {
			return "", err
		}
		userId, err := go_utils.GetLoginUserId(accessToken, code)
		if err != nil {
			return "", err
		}
		userInfo, err := go_utils.GetUserInfo(accessToken, userId)
		if err != nil {
			return "", err
		}

		sendLoginInfo(userInfo, loginCookie, accessToken)

		return userInfo.Name, nil
	}

	return "", errors.New("no login")
}

func sendLoginInfo(info *go_utils.WxUserInfo, loginCookie *CookieValue, accessToken string) {
	msg := map[string]interface{}{
		"touser": "@all", "toparty": "@all", "totag": "@all", "msgtype": "text", "agentid": conf.WxAgentId, "safe": 0,
		"text": map[string]string{
			"content": "用户[" + info.Name + "]正在扫码设置IP[" + loginCookie.OfficeIp + "]，环境[" + loginCookie.Envs + "]",
		},
	}
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken
	_, err := go_utils.HttpPost(url, msg)
	if err != nil {
		log.Println("sendLoginInfo error", err)
	}
}
