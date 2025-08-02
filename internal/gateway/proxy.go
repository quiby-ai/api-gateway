package gateway

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Route struct {
	Path        string
	Method      string
	UpstreamURL string
	RewritePath string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
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

func CORSMiddleware(corsConfig CORSConfig, allowedMethod string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if isOriginAllowed(origin, corsConfig.AllowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			if len(corsConfig.AllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
			} else {
				w.Header().Set("Access-Control-Allow-Methods", allowedMethod)
			}

			if len(corsConfig.AllowedHeaders) > 0 {
				if len(corsConfig.AllowedHeaders) == 1 && corsConfig.AllowedHeaders[0] == "*" {
					w.Header().Set("Access-Control-Allow-Headers", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
				}
			}

			if corsConfig.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != allowedMethod {
			w.Header().Set("Allow", allowedMethod)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
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
