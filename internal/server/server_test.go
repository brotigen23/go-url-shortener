package server

import (
	"log"
	"syscall"
	"testing"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRun(t *testing.T) {
	logger := zap.NewNop()
	config := &config.Config{}

	go func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	err := Run(config, logger.Sugar())
	assert.NoError(t, err)

	log.Println("DONE")
}
