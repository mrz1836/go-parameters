package parameters

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSendCORS tests the SendCORS function
func TestSendCORS(t *testing.T) {
	// Define test cases using a table-driven approach
	tests := []struct {
		name            string
		originHeader    string
		expectedHeaders map[string]string
		shouldSetOrigin bool
	}{
		{
			name:         "With Origin",
			originHeader: "https://example-origin.com",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example-origin.com",
				"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token",
				"Access-Control-Allow-Credentials": "true",
			},
			shouldSetOrigin: true,
		},
		{
			name:         "Without Origin",
			originHeader: "",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token",
				"Access-Control-Allow-Credentials": "true",
			},
			shouldSetOrigin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP request
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			if tt.originHeader != "" {
				req.Header.Set("Origin", tt.originHeader)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			// Call the SendCORS function
			SendCORS(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 OK")

			for key, expectedValue := range tt.expectedHeaders {
				value := rr.Header().Get(key)
				assert.Equal(t, expectedValue, value, "Header %s mismatch", key)
			}

			if tt.shouldSetOrigin {
				originValue := rr.Header().Get("Access-Control-Allow-Origin")
				assert.Equal(t, tt.originHeader, originValue, "Access-Control-Allow-Origin header mismatch")
			} else {
				assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"), "Access-Control-Allow-Origin should not be set")
			}
		})
	}
}

// TestFilterMap tests the FilterMap function
func TestFilterMap(t *testing.T) {
	// Set up the FilteredKeys
	FilteredKeys = []string{"password", "secret"}

	// Define test cases
	tests := []struct {
		name           string
		input          *Params
		expectedOutput *Params
	}{
		{
			name: "No keys to filter",
			input: &Params{
				Values: map[string]interface{}{
					"username": "user1",
					"email":    "user1@example.com",
				},
			},
			expectedOutput: &Params{
				Values: map[string]interface{}{
					"username": "user1",
					"email":    "user1@example.com",
				},
			},
		},
		{
			name: "Filter password",
			input: &Params{
				Values: map[string]interface{}{
					"username": "user2",
					"password": "my_password",
				},
			},
			expectedOutput: &Params{
				Values: map[string]interface{}{
					"username": "user2",
					"password": []string{FilteredValue},
				},
			},
		},
		{
			name: "Filter multiple keys",
			input: &Params{
				Values: map[string]interface{}{
					"username": "user3",
					"password": "password123",
					"secret":   "top_secret",
					"token":    "abc123",
				},
			},
			expectedOutput: &Params{
				Values: map[string]interface{}{
					"username": "user3",
					"password": []string{FilteredValue},
					"secret":   []string{FilteredValue},
					"token":    "abc123",
				},
			},
		},
		{
			name: "Value is []byte",
			input: &Params{
				Values: map[string]interface{}{
					"data": []byte("some bytes"),
				},
			},
			expectedOutput: &Params{
				Values: map[string]interface{}{
					"data": "some bytes",
				},
			},
		},
		{
			name: "Value is []byte and key is filtered",
			input: &Params{
				Values: map[string]interface{}{
					"secret": []byte("secret bytes"),
				},
			},
			expectedOutput: &Params{
				Values: map[string]interface{}{
					"secret": []string{FilteredValue},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := FilterMap(tt.input)
			assert.Equal(t, tt.expectedOutput.Values, output.Values)
		})
	}
}

// TestGeneralResponse tests the GeneralResponse function
func TestGeneralResponse(t *testing.T) {

	t.Run("Without GZIP", func(t *testing.T) {
		// Create a mock HTTP request
		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Call the GeneralResponse function
		GeneralResponse(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusOK)
		})(rr, req, nil)

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 OK")

		t.Log(rr.Header())

		assert.Empty(t, rr.Header().Get("Content-Encoding"), "Expected Content-Encoding to be empty")
	})

	t.Run("With GZIP", func(t *testing.T) {
		// Create a mock HTTP request
		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		req.Header.Set("Accept-Encoding", gZip)

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Call the GeneralResponse function
		GeneralResponse(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusOK)
		})(rr, req, nil)

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200 OK")

		t.Log(rr.Header())

		assert.Equal(t, gZip, rr.Header().Get("Content-Encoding"), "Expected Content-Encoding to be gzip")
	})
}

// TestJSONResp tests the JSONResp middleware
func TestJSONResp(t *testing.T) {
	// Track if the inner handler was called
	handlerCalled := false

	// Fake handler to wrap
	handler := func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"ok"}`))
	}

	// Wrap it with JSONResp middleware
	wrapped := JSONResp(handler)

	// Setup test server request and response
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	// httprouter requires a Params argument even if empty
	params := httprouter.Params{}
	wrapped(resp, req, params)

	// Assertions
	require.True(t, handlerCalled, "Expected handler to be called")
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "application/json", resp.Header().Get("Content-Type"))
	require.JSONEq(t, `{"message":"ok"}`, resp.Body.String())
}

// TestGeneralJSONResponse tests the GeneralJSONResponse middleware
func TestGeneralJSONResponse(t *testing.T) {
	// Track if the inner handler was called
	handlerCalled := false

	// The actual handler
	handler := func(w http.ResponseWriter, _ *http.Request) {
		handlerCalled = true
		_, _ = w.Write([]byte(`{"success":true}`))
	}

	// Wrap with GeneralJSONResponse
	wrapped := GeneralJSONResponse(handler)

	// Build request with gzip support and CORS origin
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Origin", "https://example.com")

	// Record the response
	resp := httptest.NewRecorder()

	// Call the handler
	wrapped(resp, req, httprouter.Params{})

	// Assertions
	require.True(t, handlerCalled, "Handler should have been called")
	require.Equal(t, "application/json", resp.Header().Get("Content-Type"))
	require.Equal(t, "gzip", resp.Header().Get("Content-Encoding"))
	require.Equal(t, "https://example.com", resp.Header().Get("Access-Control-Allow-Origin"))

	// Decompress the body to verify the JSON content
	reader, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()

	var body []byte
	body, err = io.ReadAll(reader)
	require.NoError(t, err)
	require.JSONEq(t, `{"success":true}`, string(body))
}
