package server

import (
	"net/http"

	"github.com/fcjr/db/internal/server/respond"
)

func (s *Server) handleGet(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	val := s.store.Get(key)

	if val == "" {
		_ = respond.Text(w,
			[]byte("key not found"),
			respond.WithStatusCode(http.StatusNotFound),
		)
		return
	}

	_ = respond.Text(w, []byte(val))
}

func (s *Server) handleSet(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	// TODO should we allow more than one set at a time?
	for key, val := range query {
		if len(val) != 1 {
			// TODO what should we do in this case?
			_ = respond.Text(w,
				[]byte("tried to set the same option more than once"),
				respond.WithStatusCode(http.StatusBadRequest),
			)
			return
		}

		s.store.Set(key, val[0])
	}

	w.WriteHeader(200)
}
