package model

type ContextData struct {
	BuildTime           string
	BundledStaticAssets []byte
	Debug               bool
	IsPartial           bool
}
