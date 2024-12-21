package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeF2le/calc_go/internal/application"
)

func TestCalcHandlerStatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{} 
		expectedStatus int          
		expectedBody   string       
	}{
		{
			name: "Valid Expression",
			body: map[string]string{
				"expression": "2+2",
			},
			expectedStatus: http.StatusOK, 
			expectedBody: `"result":4`,
		},
		{
			name: "Invalid Expression",
			body: map[string]string{
				"expression": "2+",
			},
			expectedStatus: http.StatusUnprocessableEntity, 
			expectedBody: `"error":"Expression is not valid"`,
		},
		{
			name: "Empty Expression",
			body: map[string]string{
				"expression": "",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: `"error":"Expression is not valid"`,
		},
		{
			name: "Unsupported Characters",
			body: map[string]string{
				"expression": "2+a",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: `"error":"Expression is not valid"`,
		},
		{
			name: "No Body in Request",
			body: nil,
			expectedStatus: http.StatusInternalServerError, 
			expectedBody: `"error":"Internal server error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				bodyBytes, _ := json.Marshal(tt.body)
				req = httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(bodyBytes))
			} else {
				req = httptest.NewRequest("POST", "/api/v1/calculate", nil)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			application.CalcHanlder(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("[%s] expected status %d, got %d", tt.name, tt.expectedStatus, recorder.Code)
			}

			if !bytes.Contains(recorder.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("[%s] expected body to contain '%s', got '%s'", tt.name, tt.expectedBody, recorder.Body.String())
			}
		})
	}
}
