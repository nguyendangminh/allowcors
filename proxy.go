package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type CORSReverseProxy struct {
	name   string
	prefix string
	target string
	proxy  *httputil.ReverseProxy
}

func (p *CORSReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" { // API servers might not accept the method OPTIONS
		logrus.Debugf("%s:\tAuto response to the OPTIONS method")
		return
	}
	src := r.URL.Path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, p.prefix)
	r.RequestURI = strings.TrimPrefix(r.RequestURI, p.prefix)
	logrus.Infof("%s:\t%s %s%s\t->\t%s %s%s", p.name, r.Method, r.Host, src, r.Method, p.target, r.URL.Path)
	p.proxy.ServeHTTP(w, r)
}

func NewCORSReverseProxy(name, prefix, target string) (*CORSReverseProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := NewReverseProxy(u)
	return &CORSReverseProxy{
		name:   name,
		prefix: prefix,
		target: target,
		proxy:  proxy,
	}, nil
}

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set("Host", target.Host)
		req.Host = target.Host
	}
	return &httputil.ReverseProxy{Director: director}
}

// Shameless copy from standard library
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
