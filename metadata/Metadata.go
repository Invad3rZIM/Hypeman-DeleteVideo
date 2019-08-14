package metadata

import "hypeman/comment"

//Hash = Views + Likes + Dislikes + Laughs , will be used for checking if the database data is new
type Metadata struct {
	Videoname, Artist, Songname, Bucket                   string
	Score, Comments, Views, Likes, Dislikes, Laughs, Hash int
	Tags                                                  []string
	Date                                                  int64
	CommentThread                                         []*comment.Comment
}

type DataChange struct {
	Videoname string
	Category  string
	Delta     int
}
