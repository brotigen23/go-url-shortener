package main

import (
	"fmt"
	"os"

	"github.com/brotigen23/go-url-shortener/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Printf("App run error: %v", err)
		os.Exit(1)
	}
}
