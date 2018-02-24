package main

import (
	"strings"
	"encoding/base64"
	"bytes"
	"net/http"
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

				if len(pair) == 2 &&
					string(pair[0]) == auths[0] &&
					EqualAny(string(pair[1]), auths[1:]) {
					fn(w, r) // 执行被装饰的函数
					return
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

func EqualAny(passwd string, anys []string) bool {
	for _, any := range anys {
		if any == passwd {
			return true
		}
	}

	return false
}

var auths = []string{"zywgqg", "Ca1p)Whd1s", "Qfc!jZ4ygq", "?rphlLr8yz"}
