package main

import (
	"fmt"
	slog "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/VxVxN/log"

	"reverse_proxy_server/pkg/config"
)

var protocol = "http://"

func init() {
	var err error
	if err = config.InitConfig("configs/reverse_proxy_server.yml"); err != nil {
		slog.Fatalf("Failed to init config: %v", err)
	}
	if err = log.Init("logs/reverse_proxy_server.log", getLevelLog(config.Cfg.IsTrace), false); err != nil {
		slog.Fatalf("Failed to init log: %v", err)
	}
}

func main() {
	addr := config.Cfg.Address
	log.Info.Printf("Start reverse proxy server, address %s", addr)

	http.HandleFunc("/", handleRequestAndRedirect)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal.Printf("Failed to listen and serve: %v, address: %s", err, addr)
	}
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	reqURL := getProxyURL(req.URL.String())

	if err := serveReverseProxy(reqURL, res, req); err != nil {
		log.Error.Printf("Failed to serve request: %v", err)
	}
}

func getProxyURL(url string) string {
	var host string
	for _, service := range config.Cfg.Services {
		if strings.Contains(url, service.Path) {
			host = protocol + service.Address
			break
		}
	}
	log.Trace.Printf("Request %s, redirect to %s", url, host)

	return host
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) error {
	reqURL, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("can't parse url: %v, url: %s", err, target)
	}

	proxy := httputil.NewSingleHostReverseProxy(reqURL)

	req.URL.Host = reqURL.Host
	req.URL.Scheme = reqURL.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = reqURL.Host

	proxy.ServeHTTP(res, req)
	return nil
}

func getLevelLog(isTrace bool) log.LevelLog {
	if isTrace {
		return log.TraceLog
	}
	return log.CommonLog
}
