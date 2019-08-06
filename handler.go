package main

import (
	"errors"
	"fmt"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

//EnableCors enables cors (Cross Organizational Resource Sharing)
func (h *Handler) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//VerifyBody is a helper function to ensure all http requests contain the requisite fields returns error if fields missing
func (h *Handler) verifyBody(body map[string]interface{}, str ...string) error {
	for _, s := range str {
		fmt.Println(s)
		if _, ok := body[s]; !ok {
			return errors.New("error: missing field: " + s)
		}
	}
	return nil
}
