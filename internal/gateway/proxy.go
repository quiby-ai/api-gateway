package gateway

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Route struct {
	Path        string
	Method      string
	UpstreamURL string
	RewritePath string
}

func NewProxy(upstream, rewritePath string) http.Handler {
	target, err := url.Parse(upstream)
	if err != nil {
		log.Fatalf("invalid upstream URL %q: %v", upstream, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	orig := proxy.Director
	proxy.Director = func(req *http.Request) {
		orig(req)
		if rewritePath != "" {
			req.URL.Path = rewritePath
		}
		req.Host = target.Host
	}

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf("proxy error to %s: %v", upstream, err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy
}

func MethodFilter(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}
