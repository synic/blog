package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

var (
	// build time (set during build)
	BuildTime string
	// set by the build process (so has to be exported), but use `Debug` instead since
	// it's a boolean
	DebugFlag string
	Debug     bool = false
)

func init() {
	if DebugFlag == "true" {
		godotenv.Load()
		log.Println("🐝 Debugging enabled!")
		Debug = true
	}

	if BuildTime == "" {
		BuildTime = fmt.Sprint(time.Now().Unix())
		log.Printf("🚧 Build time was not set, using %s\n", BuildTime)
	}
}
