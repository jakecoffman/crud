package crud

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServeMuxAdapter(t *testing.T) {
	t.Run("validates path parameters", func(t *testing.T) {
		adapter := NewServeMuxAdapter()
		router := NewRouter("title", "1.0", adapter)
		err := router.Add(Spec{
			Method:  "GET",
			Path:    "/widgets/{id}",
			Handler: func(w http.ResponseWriter, r *http.Request) {},
			Validate: Validate{
				Path: Object(map[string]Field{
					"id": Number().Required().Max(10),
				}),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRequest("GET", "/widgets/11", nil)
		w := httptest.NewRecorder()
		adapter.Engine.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		if !strings.Contains(w.Body.String(), `"path validation failed for field id: maximum exceeded"`) {
			t.Errorf("unexpected body %q", w.Body.String())
		}
	})

	t.Run("strips values from the body", func(t *testing.T) {
		adapter := NewServeMuxAdapter()
		router := NewRouter("title", "1.0", adapter)
		err := router.Add(Spec{
			Method: "POST",
			Path:   "/widgets",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				// reflect the body back for ease of testing
				io.Copy(w, r.Body)
			},
			Validate: Validate{
				Body: Object(map[string]Field{
					"value": String().Required(),
				}),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRequest("POST", "/widgets", strings.NewReader(`{"value": "hello", "unexpected": 1}`))
		w := httptest.NewRecorder()
		adapter.Engine.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
		// unexpected property should be stripped
		if w.Body.String() != `{"value":"hello"}` {
			t.Errorf("unexpected body %q", w.Body.String())
		}
	})

	t.Run("unexpected URL parameters are stripped", func(t *testing.T) {
		adapter := NewServeMuxAdapter()
		router := NewRouter("title", "1.0", adapter)
		err := router.Add(Spec{
			Method: "GET",
			Path:   "/widgets",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				// reflect the parameters back for ease of testing
				io.WriteString(w, r.URL.RawQuery)
			},
			Validate: Validate{
				Query: Object(map[string]Field{
					"limit": Number().Required(),
				}),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRequest("GET", "/widgets?limit=1&hello=world", nil)
		w := httptest.NewRecorder()
		adapter.Engine.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
		// unexpected property should be stripped
		if w.Body.String() != `limit=1` {
			t.Errorf("unexpected body %q", w.Body.String())
		}
	})
}
