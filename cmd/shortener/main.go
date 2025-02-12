// Модуль с точкой входа
package main

import (
	"fmt"
	"os"

	"github.com/brotigen23/go-url-shortener/internal/app"
)

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Printf("App run error: %v", err)
		os.Exit(1)
	}
}
