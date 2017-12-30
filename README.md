# AllowCORS
`AllowCORS` is a multi-backend reserve proxy. It's built to help developers to bypass same-origin policy.

## Install
```
$ go get github.com/nguyendangminh/allowcors
```

## Usage
- Write a config file `config.yml` as following:
```
--- 
port: 1203
backends: # served in order
  - 
    name: "API 1"
    path: "/api/1"
    target: "https://api.fpt.ai"
  - 
    name: "API 2"
    path: /api/2
    target: "http://api.openfpt.vn"
  - 
    name: "Web"
    path: "/"
    target: "http://vnexpress.net"

# -> http://localhost:1203/api/1/users -> https://api.fpt.ai/users
# -> http://localhost:1203/api/2/apps -> http://api.openfpt.vn/apps
# -> http://localhost:1203/* -> http://localhost:8080/*
```
- Run
```
$ allowcors -conf config.yml
$ allowcors --help
```
