package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type goIpAllowConfig struct {
	ContextPath         string
	Envs                []string // 环境
	ListenPort          int      // 监听端口
	UpdateFirewallShell string   // 更新防火墙IP脚本
	EncryptKey          string
	CookieName          string
	RedirectUri         string
}

func readIpAllowConfig() goIpAllowConfig {
	fpath := "go-ip-allow.toml"
	if len(os.Args) > 1 {
		fpath = os.Args[1]
	}

	ipAllowConfig := goIpAllowConfig{}
	if _, err := toml.DecodeFile(fpath, &ipAllowConfig); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}

	return ipAllowConfig
}

var conf goIpAllowConfig

func init() {
	conf = readIpAllowConfig()
}
