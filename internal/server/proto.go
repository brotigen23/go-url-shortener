package server

import (
	"database/sql"
	"net"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/database/migration"
	"github.com/brotigen23/go-url-shortener/internal/proto"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/repository/memory"
	"github.com/brotigen23/go-url-shortener/internal/repository/postgres"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunProto(config *config.Config, logger *zap.SugaredLogger) error {
	var repo repository.Repository
	const driver = "postgres"
	switch config.DatabaseDSN {
	case "":
		repo = memory.New(nil)
	default:
		db, err := sql.Open(driver, config.DatabaseDSN)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			return err
		}
		err = migration.MigratePostgresUp(db)
		if err != nil {
			return err
		}
		repo = postgres.New(db, logger)
	}
	service := service.New(config, logger, repo)
	server := proto.NewShortenerServer(service)

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.Errorln(err)
		return err
	}

	s := grpc.NewServer()

	proto.RegisterShotenerServer(s, server)

	logger.Infoln("Server gRPC is running")

	if err := s.Serve(listen); err != nil {
		logger.Errorln(err)
		return err
	}

	return nil
}
