package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type IpFileLine struct {
	Ip   string
	Day  string
	Envs string
	User string
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

func ParseAllowIpsFile() []IpFileLine {
	ipFileLines := make([]IpFileLine, 0)
	content, err := ioutil.ReadFile("AllowIps.txt")
	if err != nil {
		return ipFileLines
	}

	strContent := string(content)
	lines := strings.Split(strContent, "\n")

	for _, line := range lines {
		items := strings.SplitN(line, "|", 4)
		fileLine := IpFileLine{
			Ip:   items[0],
			Day:  items[1],
			Envs: items[2],
			User: items[3],
		}

		if fileLine.Ip == "" {
			continue
		}

		day, _ := time.Parse(time.RFC3339, fileLine.Day)
		if time.Now().Sub(day).Hours() > 24 {
			continue
		}

		ipFileLines = append(ipFileLines, fileLine)
	}

	return ipFileLines
}

func CreateAllowIpsFileLines(envs, ip, loginedUserName string) []IpFileLine {
	ipFileLines := make([]IpFileLine, 0)
	newIpFileLine := IpFileLine{
		Ip:   ip,
		Day:  time.Now().Format(time.RFC3339),
		Envs: envs,
		User: loginedUserName,
	}

	fileLines := ParseAllowIpsFile()

	for _, line := range fileLines {
		if line.Ip != ip {
			ipFileLines = append(ipFileLines, line)
		} else {
			mergedEnvs := mergeEnvs(newIpFileLine.Envs, line.Envs)
			if mergedEnvs != "" {
				newIpFileLine.Envs = mergedEnvs
			}
		}
	}

	return append([]IpFileLine{newIpFileLine}, ipFileLines...)
}

func mergeEnvs(s1 string, s2 string) string {
	if strings.Index(s1+",", s2+",") >= 0 {
		return ""
	}

	parts1 := strings.Split(s1, ",")
	originalPart1Len := len(parts1)
	parts2 := strings.Split(s2, ",")

	for _, part2 := range parts2 {
		if !FindInArray(part2, parts1) {
			parts1 = append(parts1, part2)
		}
	}

	if originalPart1Len == len(parts1) {
		return ""
	}

	return strings.Join(parts1, ",")

}

func FindInArray(s string, arr []string) bool {
	for _, item := range arr {
		if s == item {
			return true
		}
	}

	return false
}

func saveAllowIpsFile(ipFileLines []IpFileLine) {
	log.Println("saveAllowIpsFile:", ipFileLines)

	content := ""
	for _, line := range ipFileLines {
		if content != "" {
			content += "\n"
		}
		content += line.Ip + "|" + line.Day + "|" + line.Envs + "|" + line.User
	}

	ioutil.WriteFile("AllowIps.txt", []byte(content), 0644)
}
