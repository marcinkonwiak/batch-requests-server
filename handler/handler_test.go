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
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"requests": [{"id": "1", "method": "GET", "path": "/test"}]}`)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"requests": [{"id": "1", "method": "GET", "path": "/test"}]}`)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	h := &Handler{}
	err := h.batchRequestHandler(c)

	assert.Error(t, err)
}

func TestBatchRequestHandlerValidationError(t *testing.T) {
	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"requests": []}`)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	h := &Handler{}
	err := h.batchRequestHandler(c)

	assert.Error(t, err)
}
