# go-ip-allow
Add current ip to firewall.

![image](https://user-images.githubusercontent.com/1940588/31004430-fbb69098-a4b9-11e7-8cd1-0d6fda3941cc.png)

build:`env GOOS=linux GOARCH=amd64 go build -o go-ip-allow.linux.bin`<br/>
config file go-ip-allow.toml:

```toml
ContextPath = "/ipallow"
ListenPort = 8182
WxCorpId = "WxCorpId"
WxCorpSecret  = "WxCorpSecret"
WxAgentId = 1000003
Envs = [ "DEV", "TEST", "DEMO", "PRODUCT" ]
UpdateFirewallShell = "/home/ci/firewall/iphelp.sh"
CookieName = "ip-allow"
RedirectUri = "http://www.baidu.com"
EncryptKey = "EncryptKey"
```
bash scripts:
```bash
export http_proxy=http://127.0.0.1:9999
export https_proxy=http://127.0.0.1:9999
go get -v -u github.com/BurntSushi/toml
go get -v -u gopkg.in/kataras/iris.v6
```
fish scripts:
```fish
set -x http_proxy http://127.0.0.1:9999
set -x https_proxy http://127.0.0.1:9999
go get -v -u github.com/BurntSushi/toml
```
