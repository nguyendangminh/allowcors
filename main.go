package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var cfgFile string
var isVerbose bool
var cfg Config
var proxies map[string]*CORSReverseProxy

func main() {
	flag.StringVar(&cfgFile, "conf", "", "path to the config file")
	flag.BoolVar(&isVerbose, "v", false, "be verbose")
	flag.Parse()
	if cfgFile == "" {
		logrus.Error("Config file is required")
		Usage()
		os.Exit(1)
	}
	if isVerbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		logrus.Fatal("Reading config file failed. Error: ", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logrus.Fatal("Reading config file failed. Error: ", err)
	}

	proxies = make(map[string]*CORSReverseProxy)
	for _, b := range cfg.Backends {
		p, err := NewCORSReverseProxy(b.Name, b.Path, b.Target)
		if err != nil {
			logrus.Fatalf("%s:\tURL %s is not valid. Error %s\n", b.Name, b.Target, err)
		}
		logrus.Infof("%s:\t%s\t->\t%s", b.Name, b.Path, b.Target)
		proxies[b.Path] = p
	}
	http.HandleFunc("/", dispatch)

	logrus.Infof("Listen and serve at port %d", cfg.Port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}

func dispatch(w http.ResponseWriter, r *http.Request) {
	for _, b := range cfg.Backends {
		if strings.HasPrefix(r.URL.Path, b.Path) {
			p, _ := proxies[b.Path]
			p.ServeHTTP(w, r)
			return
		}
	}
	logrus.Errorf("%s is not matched any paths", r.URL.Path)
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
