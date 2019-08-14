package handler

import (
	"errors"
	"fmt"
	"hypeman/cache"
	"hypeman/datastore"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DataStore *datastore.DataStore
}

func NewHandler(client *mongo.Client) *Handler {
	return &Handler{
		DataStore: &datastore.DataStore{
			Client:    client,
			TimeCache: cache.NewTimeCache(client),
		},
	}
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
