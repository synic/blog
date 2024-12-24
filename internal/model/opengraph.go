package model

type OpenGraphData struct {
	Title       string `json:"title,omitempty"       yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Image       string `json:"image,omitempty"       yaml:"image,omitempty"`
	Type        string `json:"type,omitempty"        yaml:"type,omitempty"`
}
