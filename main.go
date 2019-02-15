package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func serveReverseProxy(newHost string, path string, w http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(newHost)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.URL.Path = path
	req.Host = url.Host
	req.RequestURI = ""
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	proxy.ServeHTTP(w, req)
}

func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

	url := fmt.Sprintf("http://%s", r.Host)
	newHost := strings.Replace(url, strconv.Itoa(listeningPort), strconv.Itoa(forwardingPort), 1)
	path := fmt.Sprintf("%s", r.URL.Path)

	serveReverseProxy(newHost, path, w, r)
}

var listeningPort int
var forwardingPort int

func main() {
	flag.IntVar(&listeningPort, "listeningPort", 8032, "The port to open")
	flag.IntVar(&forwardingPort, "forwardingPort", 8088, "The port to forward")
	flag.Parse()
	fmt.Println("listeningPort:  ", listeningPort)
	fmt.Println("forwardingPort: ", forwardingPort)

	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":"+strconv.Itoa(listeningPort), nil); err != nil {
		panic(err)
	}
}
