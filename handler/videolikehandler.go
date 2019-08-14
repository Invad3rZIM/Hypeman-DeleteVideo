package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) VideoLikeHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)
	var requestBody map[string]interface{}

	//ensure json is decoded
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//ensure all requisite json components are found
	if err := h.verifyBody(requestBody, "videoname", "username", "like"); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//extract field values to variables for readability
	videoname := requestBody["videoname"].(string)
	username := requestBody["username"].(string)
	like := int(requestBody["like"].(float64))

	h.DataStore.PushLikeToDB(videoname, username, like == 1)
	//ticket := opinions.NewVideoTicket(videoname, username, like, "")

	//h.DataStore.Reactions.VideoPush(ticket)

	//json.NewEncoder(w).Encode(ticket)
	w.WriteHeader(http.StatusOK)

	return
}
