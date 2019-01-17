package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)


/* This allows to modifiy the proxied response before sending it to the requester
type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("some text o patter to replace"), []byte("the replacement"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}
*/

func serveReverseProxy(newHost string, path string, w http.ResponseWriter, req *http.Request) {

	url, _ := url.Parse(newHost)

	proxy := httputil.NewSingleHostReverseProxy(url)
	// proxy.Transport = &transport{http.DefaultTransport}

	req.URL.Host   = url.Host
	req.URL.Scheme = url.Scheme
	req.URL.Path   = path
	req.Host       = url.Host
	req.RequestURI = ""
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	proxy.ServeHTTP(w, req)
}

func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

	url     := fmt.Sprintf("http://%s", r.Host)
	newHost := strings.Replace(url, listeningPort, forwardingPort, 1)
	path    := fmt.Sprintf("%s", r.URL.Path)

	serveReverseProxy(newHost, path, w, r)
}

var listeningPort = ":8032"
var forwardingPort = ":8088"

func main() {
	listeningPort = *flag.String("listeningPort", ":8032", "The port to open")
	forwardingPort   = *flag.String("forwardingPort", ":8088", "The port to forward")

	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(listeningPort, nil); err != nil {
		panic(err)
	}
}
