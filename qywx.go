package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ExpiredStore struct {
	Value   string
	Expired string
}

var (
	accessToken            string
	accessTokenExpiredTime time.Time
	accessTokenMutex       sync.Mutex
)

type TokenResult struct {
	ErrCode          int    `json:"errcode"`
	ErrMsg           string `json:"errmsg"`
	AccessToken      string `json:"access_token"`
	ExpiresInSeconds int    `json:"expires_in"`
}

func getAccessToken(corpId, corpSecret string) (string, error) {
	accessTokenMutex.Lock()
	defer accessTokenMutex.Unlock()
	if accessToken != "" && accessTokenExpiredTime.After(time.Now()) {
		return accessToken, nil
	}

	accessToken = ""

	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpId + "&corpsecret=" + corpSecret
	body, err := httpGet(url)
	if err != nil {
		return "", err
	}

	var tokenResult TokenResult
	json.Unmarshal(body, &tokenResult)
	if tokenResult.ErrCode != 0 {
		return "", errors.New(tokenResult.ErrMsg)
	}

	accessToken = tokenResult.AccessToken
	accessTokenExpiredTime = time.Now().Add(time.Duration(tokenResult.ExpiresInSeconds) * time.Second)

	return accessToken, nil
}

func createWxQyLoginUrl(redirectUri, csrfToken string) string {
	return "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=" +
		conf.WxCorpId + "&agentid=" + strconv.FormatInt(conf.WxAgentId, 10) + "&redirect_uri=" + redirectUri + "&state=" + csrfToken
}

type CookieValue struct {
	OfficeIp    string
	Envs        string
	CsrfToken   string
	ExpiredTime string
}

func writeCsrfTokenCookie(w http.ResponseWriter, csrfToken, officeIp, envs string) {
	cookieVal, err := json.Marshal(CookieValue{
		OfficeIp:    officeIp,
		Envs:        envs,
		CsrfToken:   csrfToken,
		ExpiredTime: time.Now().Add(time.Duration(24) * time.Hour).Format(time.RFC3339),
	})
	if err != nil {
		log.Println("json cookie error", err)
	}

	json := string(cookieVal)
	log.Println("csrf json:", json)
	cipher, err := CBCEncrypt(conf.EncryptKey, json)
	if err != nil {
		log.Println("CBCEncrypt cookie error", err)
	}

	cookie := http.Cookie{Name: conf.CookieName, Value: cipher, Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
}

func clearCookie(w http.ResponseWriter) {
	cookie := http.Cookie{Name: conf.CookieName, Value: "", Path: "/", Expires: time.Now().AddDate(-1, 0, 0)}
	http.SetCookie(w, &cookie)
}

func readLoginCookie(r *http.Request) *CookieValue {
	cookie, _ := r.Cookie(conf.CookieName)
	if cookie == nil {
		return nil
	}

	log.Println("cookie value:", cookie.Value)
	decrypted, _ := CBCDecrypt(conf.EncryptKey, cookie.Value)
	if decrypted == "" {
		return nil
	}

	var cookieValue CookieValue
	err := json.Unmarshal([]byte(decrypted), &cookieValue)
	if err != nil {
		log.Println("unamrshal error:", err)
		return nil
	}

	log.Println("cookie parsed:", cookieValue, ",ExpiredTime:", cookieValue.ExpiredTime)

	expired, err := time.Parse(time.RFC3339, cookieValue.ExpiredTime)
	if err != nil {
		log.Println("time.Parse:", err)
	}
	if err != nil || expired.Before(time.Now()) {
		return nil
	}

	return &cookieValue
}

type WxLoginUserId struct {
	UserId  string `json:"UserId"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func getLoginUserId(accessToken, code string) (string, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + accessToken + "&code=" + code
	body, err := httpGet(url)
	if err != nil {
		return "", err
	}

	var wxLoginUserId WxLoginUserId
	err = json.Unmarshal(body, &wxLoginUserId)
	if err != nil {
		return "", err
	}
	if wxLoginUserId.UserId == "" {
		return "", errors.New(string(body))
	}

	return wxLoginUserId.UserId, nil
}

type WxUserInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	UserId string `json:"userid"`
}

func getUserInfo(accessToken, userId string) (*WxUserInfo, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=" + accessToken + "&userid=" + userId
	body, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	var wxUserInfo WxUserInfo
	err = json.Unmarshal(body, &wxUserInfo)
	if err != nil {
		return nil, err
	}

	return &wxUserInfo, nil
}
