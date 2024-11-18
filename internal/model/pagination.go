package model

type PageData struct {
	Search     string
	Tag        string
	Items      []*Article
	TotalPages int
	Page       int
	PerPage    int
}
