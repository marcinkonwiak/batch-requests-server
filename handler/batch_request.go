package handler

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) batchRequestHandler(c echo.Context) error {
	b := new(batchRequest)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(b); err != nil {
		return err
	}

	client := &http.Client{}
	responses := make([]response, 0, len(b.Requests))

	for _, r := range b.Requests {
		resp, err := makeRequest(client, r)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		responses = append(responses, *resp)
	}

	br := newBatchResponse(&responses)
	return c.JSON(200, br)
}

func makeRequest(c *http.Client, r request) (*response, error) {
	var body bytes.Buffer
	if r.Body != nil {
		if err := json.NewEncoder(&body).Encode(r.Body); err != nil {
			return nil, err
		}
	}
	baseUrl := "https://662d27620547cdcde9e01582.mockapi.io"

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
