package app

import (
	"fmt"
	"os/exec"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
)

func Run() error {
	config := config.NewConfig()
	if config.DatabaseDSN != "" {
		cmd := exec.Command("make", "DATABASE_DSN="+config.DatabaseDSN, "migrate")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// Print the output
		fmt.Println(string(stdout))
	}

	err := server.Run(config)

	if err != nil {
		return err
	}
	return nil
}
