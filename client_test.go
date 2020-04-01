package restuss

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPerformCallAndReadResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/SUCCESS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.URL.Path == "/FAIL" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
		}))
	defer ts.Close()

	tests := []struct {
		name string
		auth *BasicAuthProvider
		path string
		want error
	}{
		{
			name: "succes",
			auth: NewBasicAuthProvider("admin", "123"),
			path: "/SUCCESS",
			want: nil,
		},
		{
			name: "fail",
			auth: NewBasicAuthProvider("admin", "123"),
			path: "/FAIL",
			want: fmt.Errorf("Retry limit exceeded"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewClient(tc.auth, ts.URL+tc.path, true)
			if err != nil {
				t.Fatalf("Error creating new client: %v", err)
			}

			req, err := http.NewRequest(http.MethodGet, c.url, nil)
			if err != nil {
				t.Fatalf("Error creating new request: %v", err)
			}

			err = c.performCallAndReadResponse(req, nil)
			if err2str(err) != err2str(tc.want) {
				t.Fatalf("got: %v, expected: %v", err, tc.want)
			}
		})
	}

}

func err2str(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
