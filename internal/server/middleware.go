package server

import (
	"fmt"
	"net/http"
)

// middleware outputs the method and uri the request is hitting on the server
func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("-> %s | %s \n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setCommon(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}

// securityHeaders wraps an http.Handler and sets a set of common security-related
// response headers to help mitigate clickjacking, MIME-sniffing, and some XSS
// risks, and to restrict referrer information.
//
// The middleware sets:
//   - X-Frame-Options: deny
//   - X-Content-Type-Options: nosniff
//   - X-XSS-Protection: 1; mode=block
//   - Referrer-Policy: strict-origin-when-cross-origin
//
// Headers are written before delegating to the wrapped handler. Use this
// middleware so these headers are present on responses from downstream handlers.
// This does not provide a Content-Security-Policy or HSTS; add those headers
// separately where required.
func (s *Server) setSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) logHandler(handlerName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("handling request",
			"method", r.Method,
			"url", r.URL.Path,
			"handler", handlerName)
		next.ServeHTTP(w, r)
	})
}
