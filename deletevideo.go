package main

import (
	"context"
	"encoding/json"
	"errors"
	"hypeman-deletevideo/handler"
)

type Data struct {
	UserId   int64  `json:"userid"`
	Password string `json:"password"`
	VideoId  string `json:"videoid"`
}

/*	Deletes the video from the S3 bucket, all Video Comments, All Video Ratings,
	And all instances of the video from the redis cache as well...

	First we must retrieve the metadata...
*/

var h *handler.Handler

func DeleteVideoHandler(ctx context.Context, data Data) (string, error) {
	//ensure all requisite json components are found
	if data.UserId == 0 {
		return "", errors.New("No userid provided")
	}
	if data.Password == "" {
		return "", errors.New("No password provided")
	}
	if data.VideoId == "" {
		return "", errors.New("No videoid provided")
	}

	var err error

	if h == nil {
		h, err = handler.InitializeHandler()

		if err != nil {
			return "", err
		}
	}

	//ensure valid user/password credentials
	err = h.Database.ValidateUserIdAndPassword(data.UserId, data.Password)

	if err != nil {
		return "", err
	}

	response, err := h.Database.DeleteVideo(data.UserId, data.VideoId)

	if err != nil {
		return "", err
	}

	//jsonify it
	j, err := json.Marshal(response)

	if err != nil {
		return "", err
	}

	return string(j), nil
}
