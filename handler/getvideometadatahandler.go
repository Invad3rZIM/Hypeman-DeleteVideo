package handler

import (
	"encoding/json"
	"fmt"
	"hypeman/comment"
	"hypeman/metadata"
	"net/http"
)

func (h *Handler) GetVideoMetadataHandler(w http.ResponseWriter, r *http.Request) {
	h.enableCors(&w)
	var requestBody map[string]interface{}

	//ensure json is decoded
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//ensure all requisite json components are found
	if err := h.verifyBody(requestBody, "videoname"); err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//extract field values to variables for readability
	videoname := requestBody["videoname"].(string)

	var metadata *metadata.Metadata
	comments := []*comment.Comment{}
	var err error

	done := make(chan bool)

	//Retrieve Metadata
	go func() {
		metadata, err = h.DataStore.TimeCache.GetMetadata(videoname)

		done <- true
		if err != nil {
			return
		}
	}()

	//Retrieve comments
	go func() {
		pointer, err := h.DataStore.GetCommentsFromDB(videoname)
		if err != nil {
			done <- true
			return
		}

		if pointer != nil {
			comments = *pointer
		}

		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}

	if metadata != nil {
		metadata.CommentThread = comments
	}

	if err != nil {
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(metadata)
	return
}
