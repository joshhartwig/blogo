package server

import (
	"net/http"
)

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	return mux
}
