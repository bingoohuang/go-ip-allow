package main

import (
	"strings"
	"encoding/base64"
	"bytes"
	"net/http"
	"time"
)

func BasicAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		basicAuthPrefix := "Basic "

		// 获取 request header
		auth := r.Header.Get("Authorization")
		// 如果是 http basic auth
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// 解码认证信息
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)

				if len(pair) == 2 {
					user := string(pair[0])
					pass := string(pair[1])

					poems := ParsePoems("./poems.txt")
					now := time.Now()
					poemsIndex := now.Day() % len(poems)
					poem := poems[poemsIndex]
					linesIndex := int(now.Weekday()) % len(poem.LinesCode)

					if user == poem.TitleCode && pass == poem.LinesCode[linesIndex] {
						fn(w, r) // 执行被装饰的函数
						return
					}
				}
			}
		}
		w.Header().Set("Content-Type", "'Content-type:text/html;charset=ISO-8859-1'")
		// 认证失败，提示 401 Unauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401 状态码
		w.WriteHeader(http.StatusUnauthorized)
	}
}
