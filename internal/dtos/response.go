package dtos

type Response struct {
	RequestID string              `json:"request_id,omitempty"`
	Data      any                 `json:"data,omitempty"`
	Errors    map[string][]string `json:"errors,omitempty"`
}
