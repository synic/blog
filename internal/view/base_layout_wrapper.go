package view

import (
	"github.com/a-h/templ"

	"github.com/synic/adamthings.me/internal/model"
)

type baseLayoutConfig struct {
	title string
	og    model.OpenGraphData
}

func getDefaultBaseLayoutConfig() baseLayoutConfig {
	return baseLayoutConfig{}
}

func WithOpenGraphData(og model.OpenGraphData) func(*baseLayoutConfig) {
	return func(conf *baseLayoutConfig) {
		conf.og = og
	}
}

func BaseLayout(title string, options ...func(*baseLayoutConfig)) templ.Component {
	conf := getDefaultBaseLayoutConfig()

	for _, option := range options {
		option(&conf)
	}

	return baseLayout(conf)
}
