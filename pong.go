package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PongResponse represents the HTTP response of the PongHandler func. Its purpose is to store an
// incoming request's HTTP method, query params and body.
type PongResponse struct {
	Method      string      `json:"method"`
	QueryParams interface{} `json:"queryParams"`
	Body        interface{} `json:"body"`
}

// PongHandler is an HTTP handleFunc that renders a PongResponse, thus displaying the incoming
// request's method, query params and body.
func PongHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pongResponse := PongResponse{
		Method:      r.Method,
		QueryParams: r.URL.Query(),
		Body:        r.PostForm,
	}

	if response, err := json.Marshal(&pongResponse); err != nil {
		errorText := fmt.Errorf("json.Marshal(): %s", err).Error()
		http.Error(w, errorText, http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(response))
	}
}
