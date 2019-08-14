package handler

import (
	"encoding/json"
	"fmt"
	"hypeman/metadata"
	"net/http"
)

func (h *Handler) GlobalMostXHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)
	var requestBody map[string]interface{}

	//ensure json is decoded
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//ensure all requisite json components are found
	if err := h.verifyBody(requestBody, "category", "needed", "startindex", "time"); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//extract field values to variables for readability
	category := requestBody["category"].(string)
	needed := int(requestBody["needed"].(float64))
	startIndex := int(requestBody["startindex"].(float64))
	ti := requestBody["time"].(string)

	var mds *[]*metadata.Metadata

	mds, err := h.DataStore.TimeCache.GetVideos(ti, category, needed, startIndex)

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(&mds)
	w.WriteHeader(http.StatusOK)

	return
}
