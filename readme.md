**Build**

`$go mod tidy`<br>
`$go build -o testweb main.go`

**Deployment**

* testweb -> /usr/local/bin/testweb
* nginx.conf -> /etc/nginx/nginx.conf
* override.conf -> /etc/systemd/system/testweb.service.d/override.conf
* testweb.conf -> /etc/nginx/conf.d/testweb.conf
* testweb.service -> /usr/lib/systemd/system/testweb.service