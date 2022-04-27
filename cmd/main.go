package main

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	repo "github.com/alexzh7/sample-service/internal/dvdstore/repository/postgres"
	"github.com/alexzh7/sample-service/internal/dvdstore/usecase"
	"github.com/alexzh7/sample-service/internal/models"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

func main() {
	// Logger
	l := zap.NewExample().Sugar()
	defer l.Sync()

	// TODO: Read config from env/file
	// TODO: Fix TODOs :)
	connStr := "user=pguser password=pgpass dbname=dvdstore sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal(err)
	}

	pgRepo, err := repo.NewPgRepo(db)
	if err != nil {
		l.Fatal(err)
	}

	validator := validator.New()

	uc := usecase.NewDvdstoreUC(pgRepo, l, validator)

	prods := []*models.Product{{Quantity: 2}}
	res, err := uc.AddOrder(10, prods)

	fmt.Println(res)
	fmt.Println(err)
}
