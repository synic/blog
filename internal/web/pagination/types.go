package pagination

import (
	"github.com/synic/adamthings.me/internal/model"
)

type PageData struct {
	Search     string
	Tag        string
	Items      []*model.Article
	TotalPages int
	Page       int
	PerPage    int
}
