package comment

import (
	"fmt"
	"time"
)

type Comment struct {
	Videoname string
	CommentID string
	Date      int64
	Username  string
	Body      string
	Score     int
}

//NewComment returns a timestamped wrapper for username + string
func NewComment(thread string, username string, body string) *Comment {
	date := time.Now().Unix()

	return &Comment{
		Username:  username,
		Body:      body,
		Score:     0,
		Videoname: thread,
		CommentID: fmt.Sprintf("%d*%s", date, username),
		Date:      date,
	}
}
