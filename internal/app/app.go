// Модуль app обеспечивает инициализацию необходимых модулей для сервера и последующего его запуска.
package app

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLog() (*zap.SugaredLogger, error) {
	// Init
	// Logger
	logConf := zap.NewDevelopmentConfig()
	logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logConf.ErrorOutputPaths = []string{"./errors.log"}
	logger, err := logConf.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func Run() error {

	logger, err := initLog()
	if err != nil {
		return err
	}

	err = godotenv.Load()
	if err != nil {
		logger.Debugln(err.Error())
	}
	config, err := config.NewConfig()
	if err != nil {
		return err
	}
	logger.Infoln("config is loaded", config)
	err = server.Run(config, logger)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	return nil
}
