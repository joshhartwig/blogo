package main

import (
	"fmt"
	"net/http"
)

func (c *config) blogHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *config) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (c *config) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(c.tplCache)
	err := c.render(w, "home", nil)
	if err != nil {
		c.logger.Error("error rendering template", err.Error(), "error")
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}

// middleware outputs the method and uri the request is hitting on the server
func (c *config) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("<= %s | %s \n", r.Method, r.URL)
		w.Header().Set("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}
