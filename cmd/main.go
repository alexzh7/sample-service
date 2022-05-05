package main

import (
	"go.uber.org/zap"

	"github.com/alexzh7/sample-service/config"
	"github.com/alexzh7/sample-service/internal/server"
	"github.com/alexzh7/sample-service/pkg/postgres"
	_ "github.com/lib/pq"
)

func main() {
	// Create logger
	l := zap.NewExample().Sugar()
	defer l.Sync()

	// Load config
	config, err := config.NewConfig()
	if err != nil {
		l.Fatalf("Config init: %v", err)
	}

	// Create db connection
	dbConn, err := postgres.NewPostgresConn(config)
	if err != nil {
		l.Fatalf("Postgresql init: %v", err)
	}

	// Run server
	s := server.NewServer(config, l, dbConn)
	if err := s.Run(); err != nil {
		l.Fatalf("DVDStore server init: %v", err)
	}
}
