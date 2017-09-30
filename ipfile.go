package main

import (
	"io/ioutil"
	"strings"
	"time"
)

type IpFileLine struct {
	Ip  string
	Day string
}

func joinAllowedIpLines(lines []IpFileLine) string {
	content := ""
	for _, line := range lines {
		if content != "" {
			content += ","
		}
		content += line.Ip
	}
	return content
}

func parseAllowIpsFile(envs, ip string) []IpFileLine {
	ipFileLines := make([]IpFileLine, 0)
	ipFileLines = append(ipFileLines, IpFileLine{
		Ip:  ip,
		Day: time.Now().Format(time.RFC3339),
	})

	content, err := ioutil.ReadFile(envs + "-AllowIps.txt")
	if err != nil {
		return ipFileLines
	}

	strContent := string(content)
	lines := strings.Split(strContent, "\n")

	for _, line := range lines {
		items := strings.SplitN(line, " ", 2)
		fileLine := IpFileLine{
			Ip:  items[0],
			Day: items[1],
		}

		day, _ := time.Parse(time.RFC3339, fileLine.Day)
		if time.Now().Sub(day).Hours() > 24 {
			continue
		}

		if fileLine.Ip != ip {
			ipFileLines = append(ipFileLines, fileLine)
		}
	}

	return ipFileLines
}

func saveAllowIpsFile(env string, ipFileLines []IpFileLine) {
	content := ""
	for _, line := range ipFileLines {
		if content != "" {
			content += "\n"
		}
		content += line.Ip + " " + line.Day
	}

	ioutil.WriteFile(env+"-AllowIps.txt", []byte(content), 0644)
}
