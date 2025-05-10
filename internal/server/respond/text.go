package respond

import "net/http"

type respondOptions struct {
	statusCode int
}

func WithStatusCode(statusCode int) func(opts *respondOptions) {
	return func(opts *respondOptions) {
		opts.statusCode = statusCode
	}
}

func Text(w http.ResponseWriter, text []byte, opts ...func(*respondOptions)) error {
	options := &respondOptions{
		statusCode: http.StatusOK,
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
