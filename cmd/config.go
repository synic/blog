package cmd

// set to debug from debug.go if `debug` build tag is used
var isDebugBuild = false

type config struct {
	Bind string `env:"BIND,default=:3000"`
}
