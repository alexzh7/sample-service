package main

import (
	"database/sql"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	service "github.com/alexzh7/sample-service/internal/dvdstore/grpc"
	repo "github.com/alexzh7/sample-service/internal/dvdstore/repository/postgres"
	"github.com/alexzh7/sample-service/internal/dvdstore/usecase"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/alexzh7/sample-service/proto"
	_ "github.com/lib/pq"
)

const grpcport = 9090

// TODO: Fix TODOs :)
func main() {
	// Logger
	l := zap.NewExample().Sugar()
	defer l.Sync()

	// TODO: Read config from env/file
	connStr := "user=pguser password=pgpass dbname=dvdstore sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal(err)
	}

	// New repo
	pgRepo, err := repo.NewPgRepo(db)
	if err != nil {
		l.Fatal(err)
	}
	validator := models.NewValidation()
	// New use case
	uc := usecase.NewDvdstoreUC(pgRepo, l, validator)

	// New grpc server
	grpcSrv := grpc.NewServer()
	grpcService := service.NewDvdstoreService(uc, l)
	proto.RegisterDvdstoreServer(grpcSrv, grpcService)

	reflection.Register(grpcSrv)

	// Create a TCP socket for inbound server connections
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcport))
	if err != nil {
		l.Fatalf("Unable to create listener: %v", err)
	}
	l.Infof("GRPC listening on port %v", grpcport)

	//TODO: graceful shutdown
	grpcSrv.Serve(ls)

	// grpcurl -plaintext localhost:9090 describe
	// grpcurl -plaintext localhost:9090 describe proto.Dvdstore.GetOrder
	// grpcurl -plaintext localhost:9090 describe proto.GetOrderReq
	// grpcurl -d '{"Limit":20}' -plaintext localhost:9090 proto.Dvdstore/GetCustomers
}
