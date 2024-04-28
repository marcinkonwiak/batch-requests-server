package handler

type batchRequest struct {
	Requests []request `json:"requests" validate:"min=1,max=20,required,dive"`
}

type request struct {
	Id      string            `json:"id" validate:"required"`
	Path    string            `json:"path" validate:"required,relativeUrl,allowedPath"`
	Method  string            `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE HEAD OPTIONS TRACE CONNECT"`
	Body    interface{}       `json:"body,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}
