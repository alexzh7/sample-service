package server

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/alexzh7/sample-service/config"
	service "github.com/alexzh7/sample-service/internal/dvdstore/grpc"
	repo "github.com/alexzh7/sample-service/internal/dvdstore/repository/postgres"
	"github.com/alexzh7/sample-service/internal/dvdstore/usecase"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/alexzh7/sample-service/proto"
	"go.uber.org/zap"
)

// Server is application server struct
type Server struct {
	config *config.Config
	log    *zap.SugaredLogger
	dbConn *sql.DB
}

// NewServer returns new application server
func NewServer(config *config.Config, log *zap.SugaredLogger, dbConn *sql.DB) *Server {
	return &Server{config: config, log: log, dbConn: dbConn}
}

func (s *Server) Run() error {
	// New repository
	pgRepo, err := repo.NewPgRepo(s.dbConn)
	if err != nil {
		s.log.Fatal(err)
	}

	// New validator
	validator := models.NewValidation()

	// New use case
	uc := usecase.NewDvdstoreUC(pgRepo, s.log, validator)

	// New grpc server
	grpcSrv := grpc.NewServer()
	grpcService := service.NewDvdstoreService(uc, s.log)
	proto.RegisterDvdstoreServer(grpcSrv, grpcService)

	reflection.Register(grpcSrv)

	// Create a TCP socket for inbound server connections
	grpcport := s.config.GRPC.Port
	ls, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcport))
	if err != nil {
		s.log.Fatalf("Unable to create listener: %v", err)
	}

	// Start GRPC server and shutdown gracefully
	go func() {
		s.log.Infof("GRPC listening on port %v", grpcport)
		s.log.Fatal(grpcSrv.Serve(ls))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	grpcSrv.GracefulStop()
	s.log.Info("Server exited properly")

	return nil
}
