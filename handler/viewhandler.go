package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) VideoViewHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)
	var requestBody map[string]interface{}

	//ensure json is decoded
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//ensure all requisite json components are found
	if err := h.verifyBody(requestBody, "videoname", "username"); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//extract field values to variables for readability
	videoname := requestBody["videoname"].(string)
	username := requestBody["username"].(string)

	h.DataStore.PushViewToDB(videoname, username)

	json.NewEncoder(w).Encode("ok")
	w.WriteHeader(http.StatusOK)

	return
}
