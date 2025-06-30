package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ForwardTo(baseURL string) http.HandlerFunc {
	target, _ := url.Parse(baseURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(w http.ResponseWriter, r *http.Request) {

		proxy.ServeHTTP(w, r)
		// take any http request, send it to the real server, return the response back to the user
	}
}
