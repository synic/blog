package model

import "time"

type Comment struct {
	ID        int64
	Body      string
	Username  string
	AvatarURL string
	CreatedAt time.Time
	ParentID  *int64
}

type CommentThread struct {
	Comment
	Replies []Comment
}

func OrganizeComments(comments []Comment) []CommentThread {
	threads := []CommentThread{}
	replyMap := map[int64][]Comment{}

	for _, c := range comments {
		if c.ParentID == nil {
			threads = append(threads, CommentThread{Comment: c})
		} else {
			replyMap[*c.ParentID] = append(replyMap[*c.ParentID], c)
		}
	}

	for i := range threads {
		threads[i].Replies = replyMap[threads[i].ID]
	}

	return threads
}
