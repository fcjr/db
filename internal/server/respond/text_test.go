package respond_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fcjr/db/internal/server/respond"
	"github.com/fcjr/db/internal/utils"
)

func TestText(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name       string
		response   string
		statusCode *int
		expectErr  bool
	}{
		{
			name:       "Text should return the text provided, with a 200 status code by default",
			response:   "test",
			statusCode: nil,
		},
		{
			name:       "Text should return the text provided, with the status code provided",
			response:   "test",
			statusCode: utils.ToPtr(201),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			wr := httptest.NewRecorder()
			var err error
			if tc.statusCode != nil {
				err = respond.Text(wr, []byte(tc.response), respond.WithStatusCode(*tc.statusCode))
			} else {
				err = respond.Text(wr, []byte(tc.response))
			}

			// check error
			if tc.expectErr && err != nil {
				t.Fatal(err)
			}

			// check body
			if wr.Body.String() != tc.response {
				t.Fatalf("expected %s for body, got %s", tc.response, wr.Body.String())
			}

			// check content type
			expectedType := "text/plain; charset=utf-8"
			if wr.Header().Get("Content-Type") != expectedType {
				t.Fatalf("expected %s for content type, got %s", expectedType, wr.Header().Get("Content-Type"))
			}

			// check status code
			resultStatus := wr.Result().StatusCode
			if tc.statusCode != nil {
				if resultStatus != *tc.statusCode {
					t.Fatalf("expected %d for status code, got %d", *tc.statusCode, resultStatus)
				}
			} else if resultStatus != http.StatusOK {
				t.Fatalf("expected %d for default status code, got %d", http.StatusOK, resultStatus)
			}

		})
	}
}
