package server

import "net/http"

func (s *Server) handleGet(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	val := s.store.Get(key)

	_ = respondText(w, []byte(val))
}

func (s *Server) handleSet(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	// TODO should we allow more than one set at a time?
	for key, val := range query {
		if len(val) != 1 {
			// TODO what should we do in this case?
			_ = respondText(w,
				[]byte("tried to set the same option more than once"),
				withStatusCode(http.StatusBadRequest),
			)
			return
		}

		s.store.Set(key, val[0])
	}

	w.WriteHeader(200)
}
