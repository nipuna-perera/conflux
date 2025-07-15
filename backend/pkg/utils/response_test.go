package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantStatus int
		wantData   interface{}
	}{
		{
			name:       "success response with object",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
			wantData:   map[string]interface{}{"message": "success"},
		},
		{
			name:       "created response",
			statusCode: http.StatusCreated,
			data:       map[string]interface{}{"id": 123, "name": "test"},
			wantStatus: http.StatusCreated,
			wantData:   map[string]interface{}{"id": float64(123), "name": "test"}, // JSON unmarshals numbers as float64
		},
		{
			name:       "error response",
			statusCode: http.StatusBadRequest,
			data:       map[string]string{"error": "bad request"},
			wantStatus: http.StatusBadRequest,
			wantData:   map[string]interface{}{"error": "bad request"},
		},
		{
			name:       "nil data",
			statusCode: http.StatusNoContent,
			data:       nil,
			wantStatus: http.StatusNoContent,
			wantData:   nil,
		},
		{
			name:       "empty string",
			statusCode: http.StatusOK,
			data:       "",
			wantStatus: http.StatusOK,
			wantData:   "",
		},
		{
			name:       "array data",
			statusCode: http.StatusOK,
			data:       []string{"item1", "item2"},
			wantStatus: http.StatusOK,
			wantData:   []interface{}{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test response recorder
			w := httptest.NewRecorder()

			// Call function
			JSONResponse(w, tt.statusCode, tt.data)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("JSONResponse() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("JSONResponse() Content-Type = %v, want application/json", contentType)
			}

			// Check response body if data is not nil
			if tt.data != nil {
				var got interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Errorf("JSONResponse() failed to unmarshal response: %v", err)
					return
				}

				// Compare data
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.wantData)
				if !bytes.Equal(gotJSON, wantJSON) {
					t.Errorf("JSONResponse() body = %v, want %v", string(gotJSON), string(wantJSON))
				}
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		wantStatus int
	}{
		{
			name:       "bad request error",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unauthorized error",
			statusCode: http.StatusUnauthorized,
			message:    "Authentication required",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "not found error",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			message:    "Something went wrong",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "empty message",
			statusCode: http.StatusBadRequest,
			message:    "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test response recorder
			w := httptest.NewRecorder()

			// Call function
			ErrorResponse(w, tt.statusCode, tt.message)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("ErrorResponse() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("ErrorResponse() Content-Type = %v, want application/json", contentType)
			}

			// Check response structure
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("ErrorResponse() failed to unmarshal response: %v", err)
				return
			}

			// Verify error response structure
			if response["error"] != true {
				t.Errorf("ErrorResponse() error field = %v, want true", response["error"])
			}

			if response["message"] != tt.message {
				t.Errorf("ErrorResponse() message = %v, want %v", response["message"], tt.message)
			}

			if response["status"] != float64(tt.statusCode) { // JSON unmarshals numbers as float64
				t.Errorf("ErrorResponse() status = %v, want %v", response["status"], tt.statusCode)
			}
		})
	}
}

func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantData interface{}
	}{
		{
			name:     "string data",
			data:     "success message",
			wantData: "success message",
		},
		{
			name:     "object data",
			data:     map[string]interface{}{"id": 1, "name": "test"},
			wantData: map[string]interface{}{"id": float64(1), "name": "test"},
		},
		{
			name:     "array data",
			data:     []string{"item1", "item2"},
			wantData: []interface{}{"item1", "item2"},
		},
		{
			name:     "nil data",
			data:     nil,
			wantData: nil,
		},
		{
			name:     "empty object",
			data:     map[string]interface{}{},
			wantData: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test response recorder
			w := httptest.NewRecorder()

			// Call function
			SuccessResponse(w, tt.data)

			// Check status code
			if w.Code != http.StatusOK {
				t.Errorf("SuccessResponse() status = %v, want %v", w.Code, http.StatusOK)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("SuccessResponse() Content-Type = %v, want application/json", contentType)
			}

			// Check response structure
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("SuccessResponse() failed to unmarshal response: %v", err)
				return
			}

			// Verify success response structure
			if response["error"] != false {
				t.Errorf("SuccessResponse() error field = %v, want false", response["error"])
			}

			// Compare data
			gotData := response["data"]
			gotJSON, _ := json.Marshal(gotData)
			wantJSON, _ := json.Marshal(tt.wantData)
			if !bytes.Equal(gotJSON, wantJSON) {
				t.Errorf("SuccessResponse() data = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

// Test edge cases and error conditions
func TestJSONResponseEdgeCases(t *testing.T) {
	t.Run("complex nested data", func(t *testing.T) {
		w := httptest.NewRecorder()
		complexData := map[string]interface{}{
			"nested": map[string]interface{}{
				"array": []interface{}{1, 2, 3, "string", map[string]string{"inner": "value"}},
				"bool":  true,
				"null":  nil,
			},
			"number": 42.5,
		}

		JSONResponse(w, http.StatusOK, complexData)

		if w.Code != http.StatusOK {
			t.Errorf("JSONResponse() status = %v, want %v", w.Code, http.StatusOK)
		}

		var got map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("JSONResponse() failed to unmarshal complex data: %v", err)
		}
	})
}

// Benchmark tests
func BenchmarkJSONResponse(b *testing.B) {
	data := map[string]interface{}{
		"id":      123,
		"name":    "benchmark test",
		"active":  true,
		"details": map[string]string{"description": "test object"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		JSONResponse(w, http.StatusOK, data)
	}
}

func BenchmarkErrorResponse(b *testing.B) {
	message := "Benchmark error message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		ErrorResponse(w, http.StatusBadRequest, message)
	}
}

func BenchmarkSuccessResponse(b *testing.B) {
	data := map[string]interface{}{
		"result": "success",
		"count":  100,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		SuccessResponse(w, data)
	}
}
