//go:build debug

package cmd

import (
	"github.com/joho/godotenv"
)

func init() {
	isDebugBuild = true
	godotenv.Overload(".env")
}
