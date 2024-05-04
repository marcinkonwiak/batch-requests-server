package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func (h *Handler) batchRequestHandler(c echo.Context) error {
	b := new(batchRequest)
	if err := c.Bind(b); err != nil {
		return h.processBindingError(err)
	}
	if err := c.Validate(b); err != nil {
		return err
	}

	client := &http.Client{}
	responses := make(chan response, len(b.Requests))
	limiter := make(chan struct{}, viper.GetInt("max_concurrent_requests"))

	for _, r := range b.Requests {
		limiter <- struct{}{}
		go func() {
			defer func() { <-limiter }()
			resp, err := makeRequest(client, r)
			if err != nil {
				resp = newErrorResponse(r.Id, err)
			}
			responses <- *resp
		}()
	}

	// Limit the number of concurrent requests
	for i := 0; i < cap(limiter); i++ {
		limiter <- struct{}{}
	}

	responsesOutput := make([]response, 0, len(b.Requests))
	for range len(b.Requests) {
		resp := <-responses
		responsesOutput = append(responsesOutput, resp)
	}

	return c.JSON(200, newBatchResponse(&responsesOutput))
}

func makeRequest(c *http.Client, r request) (*response, error) {
	var body bytes.Buffer
	if r.Body != nil {
		if err := json.NewEncoder(&body).Encode(r.Body); err != nil {
			return nil, err
		}
	}

	baseUrl := viper.GetString("base_url")
	req, err := http.NewRequest(r.Method, baseUrl+r.Path, &body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	o := newResponse(r.Id, resp)
	return o, nil
}
