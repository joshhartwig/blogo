package main

import (
	"fmt"
	"net/http"
)

func (c *config) blogHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(c.tplCache)
	err := c.tplCache["home.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (c *config) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := c.render(w, "home", nil)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}

// middleware outputs the method and uri the request is hitting on the server
func (c *config) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("<= %s | %s \n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
