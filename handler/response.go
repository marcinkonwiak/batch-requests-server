package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

type batchResponse struct {
	Responses []response `json:"responses"`
}

func newBatchResponse(r *[]response) *batchResponse {
	responses := make([]response, len(*r))
	for i, v := range *r {
		responses[i] = v
	}

	return &batchResponse{
		Responses: responses,
	}
}

type response struct {
	Id         string              `json:"id"`
	StatusCode int                 `json:"statusCode"`
	Body       interface{}         `json:"body"`
	Headers    map[string][]string `json:"headers"`
}

func newResponse(id string, r *http.Response) *response {
	var body interface{}

	b, err := io.ReadAll(r.Body)

	if err != nil {
		body = map[string]interface{}{
			"error": "Error parsing response body",
		}
	}

	if err = json.Unmarshal(b, &body); err != nil {
		body = string(b)
	}

	return &response{
		Id:         id,
		StatusCode: r.StatusCode,
		Body:       body,
		Headers:    r.Header,
	}
}

func newErrorResponse(id string, err error) *response {
	return &response{
		Id:         id,
		StatusCode: 500,
		Body: map[string]interface{}{
			"batchApiError": err.Error(),
		},
		Headers: nil,
	}
}
