package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type CORSReverseProxy struct {
	proxy *httputil.ReverseProxy
}

func (p *CORSReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	p.proxy.ServeHTTP(w, r)
}

func main() {
	var targetURL string
	var port int
	var cors bool
	flag.StringVar(&targetURL, "url", "http://minhnd.com", "Target URL")
	flag.IntVar(&port, "p", 1203, "Port to bind")
	flag.BoolVar(&cors, "cors", false, "Enable CORS")
	flag.Parse()

	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal("Parse URL failed: ", err)
	}

	var proxy http.Handler
	if cors {
		proxy = NewCORSReverseProxy(target)
	} else {
		proxy = NewReverseProxy(target)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), proxy))
}

func NewCORSReverseProxy(target *url.URL) *CORSReverseProxy {
	proxy := NewReverseProxy(target)
	return &CORSReverseProxy{proxy: proxy}
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
		log.Println(req.Host)
		log.Println(req.URL)
	}
	return &httputil.ReverseProxy{Director: director}
}

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
