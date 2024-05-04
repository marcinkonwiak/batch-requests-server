package handler

import (
	"bytes"
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/marcinkonwiak/batch-requests-server/validator"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setup() *echo.Echo {
	e := echo.New()
	e.Validator = validator.NewValidator()

	viper.Set("allowed_paths", []string{"^.*"})
	viper.Set("max_concurrent_requests", 1)
	viper.Set("base_url", "http://localhost:8080")

	return e
}

func TestBatchRequestHandlerSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:8080/test",
		httpmock.NewStringResponder(200, `{"test": "value"}`))

	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(
		`{"requests": [{"id": "1", "method": "GET", "path": "/test"}]}`,
	)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{}
	err := h.batchRequestHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res batchResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.Len(t, res.Responses, 1)
	assert.Equal(t, "1", res.Responses[0].Id)
	assert.Equal(t, 200, res.Responses[0].StatusCode)
	assert.Equal(t, map[string]interface{}{"test": "value"}, res.Responses[0].Body)
	assert.Equal(t, map[string][]string{}, res.Responses[0].Headers)
}

func TestBatchRequestHandlerConnectionRefused(t *testing.T) {
	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(
		`{"requests": [{"id": "1", "method": "GET", "path": "/test"}]}`,
	)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{}
	err := h.batchRequestHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res batchResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.Responses[0].StatusCode)
}

func TestBatchRequestHandlerBindingError(t *testing.T) {
	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"requests": "invalid"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{}
	err := h.batchRequestHandler(c)

	assert.Error(t, err)
}

func TestBatchRequestHandlerValidationError(t *testing.T) {
	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"requests": []}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{}
	err := h.batchRequestHandler(c)

	assert.Error(t, err)
}

func TestBatchRequestHandlerRelativeUrlValidation(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expectErr bool
	}{
		{
			name:      "Wrong path",
			data:      `{"requests": [{"id": "1", "method": "GET", "path": "wrong/path"}]}`,
			expectErr: true,
		},
		{
			name:      "Full URL",
			data:      `{"requests": [{"id": "1", "method": "GET", "path": "https://wrong.path/test"}]}`,
			expectErr: true,
		},
		{
			name:      "Correct path",
			data:      `{"requests": [{"id": "1", "method": "GET", "path": "/test"}]}`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setup()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.data)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &Handler{}
			err := h.batchRequestHandler(c)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBatchRequestHandlerAllowedPathValidation(t *testing.T) {
	tests := []struct {
		name        string
		allowedPath string
		expectErr   bool
	}{
		{
			name:        "Correct path",
			allowedPath: "^/test$",
			expectErr:   false,
		},
		{
			name:        "Wrong path",
			allowedPath: "^/wrong/path$",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setup()
			viper.Set("allowed_paths", []string{tt.allowedPath})
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(
				`{"requests": [{"id": "1", "method": "GET", "path": "/test"}]}`),
			))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &Handler{}
			err := h.batchRequestHandler(c)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
