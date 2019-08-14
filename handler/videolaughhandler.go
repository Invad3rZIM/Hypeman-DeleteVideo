package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) VideoLaughHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)
	var requestBody map[string]interface{}

	//ensure json is decoded
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//ensure all requisite json components are found
	if err := h.verifyBody(requestBody, "videoname", "username", "laugh"); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//extract field values to variables for readability
	videoname := requestBody["videoname"].(string)
	username := requestBody["username"].(string)
	laugh := int(requestBody["laugh"].(float64))

	h.DataStore.PushLaughToDB(videoname, username, laugh == 1)

	json.NewEncoder(w).Encode("ok")
	w.WriteHeader(http.StatusOK)

	return
}
