package server

import "net/http"

type respondTextOptions struct {
	statusCode int
}

func withStatusCode(statusCode int) func(opts *respondTextOptions) {
	return func(opts *respondTextOptions) {
		opts.statusCode = statusCode
	}
}

func respondText(w http.ResponseWriter, text []byte, opts ...func(*respondTextOptions)) error {
	options := &respondTextOptions{
		statusCode: http.StatusAccepted,
	}

	for _, opt := range opts {
		opt(options)
	}

	w.WriteHeader(options.statusCode)
	w.Header().Add("Content-Type", "application/text")

	_, err := w.Write(text)
	if err != nil {
		return err
	}

	return nil
}
