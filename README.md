# go-ip-allow
Add current ip to firewall.

![image](https://user-images.githubusercontent.com/1940588/31004430-fbb69098-a4b9-11e7-8cd1-0d6fda3941cc.png)

build:
`./gobin.sh`
`env GOOS=linux GOARCH=amd64 go build -o go-ip-allow.linux.bin`<br/>
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
```
fish scripts:
```fish
set -x http_proxy http://127.0.0.1:9999
set -x https_proxy http://127.0.0.1:9999
go get -v -u github.com/BurntSushi/toml
```


## Nginx config
```nginx
    server {
        listen 8000;
        server_name ip.a.b;

        location / {
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header Host $host;
            proxy_pass http://localhost:8182;
        }
      }
```