package internal

import (
	"fmt"
	"log"
	"time"
)

var (
	// build time (set during build)
	BuildTime string
	DebugFlag string
	Debug     bool = false
)

func init() {
	if DebugFlag == "true" {
		Debug = true
	}

	if BuildTime == "" {
		BuildTime = fmt.Sprint(time.Now().Unix())
		log.Printf("⚠️ Build time was not set, using %s\n", BuildTime)
	}
}
