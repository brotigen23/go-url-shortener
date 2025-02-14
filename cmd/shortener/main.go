// Модуль с точкой входа
package main

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func initTags() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}

func main() {
	initTags()
	err := app.Run()
	if err != nil {
		fmt.Printf("App run error: %v", err)
	}
}
