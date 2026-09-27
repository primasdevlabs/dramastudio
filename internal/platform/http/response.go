package http

import (
	"encoding/json"
	"net/http"
)

type APIErrorDetail struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type APIErrorResponse struct {
	Error APIErrorDetail `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func WriteError(w http.ResponseWriter, status int, code, message, requestID string, details interface{}) {
	if details == nil {
		details = map[string]interface{}{}
	}
	resp := APIErrorResponse{
		Error: APIErrorDetail{
			Code:      code,
			Message:   message,
			Details:   details,
			RequestID: requestID,
		},
	}
	WriteJSON(w, status, resp)
}
