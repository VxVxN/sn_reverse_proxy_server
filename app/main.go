package main

import (
	slog "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/VxVxN/log"

	"reverse_proxy_server/app/config"
)

var protocol = "http://"

func init() {
	var err error
	if err = log.Init("app/logs/reverse_proxy_server.log", log.CommonLog, false); err != nil {
		slog.Printf("Failed to init log: %v", err)
	}
	if err = config.InitConfig("app/config/main.json"); err != nil {
		log.Fatal.Printf("Failed to init config: %v", err)
	}
}

func main() {
	log.Info.Println("Reverse proxy server start.")

	http.HandleFunc("/", handleRequestAndRedirect)

	addr := config.Cfg.ReverseProxyServerHostname + ":" + strconv.Itoa(config.Cfg.ReverseProxyServerPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal.Printf("Failed to listen and serve: %v, address: %s", err, addr)
	}
}

func getProxyURL(url string) string {
	var host string
	if strings.Contains(url, "/ajax/") {
		host = protocol + config.Cfg.AJAXServerHostname + ":" + strconv.Itoa(config.Cfg.AJAXServerPort)
		return host
	}

	host = protocol + config.Cfg.WebServerHostname + ":" + strconv.Itoa(config.Cfg.WebServerPort)
	return host
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	reqURL, err := url.Parse(target)
	if err != nil {
		log.Error.Printf("Failed to parse url: %v, url: %s", err, target)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(reqURL)

	req.URL.Host = reqURL.Host
	req.URL.Scheme = reqURL.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = reqURL.Host

	proxy.ServeHTTP(res, req)
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	reqURL := getProxyURL(req.URL.String())

	serveReverseProxy(reqURL, res, req)
}
