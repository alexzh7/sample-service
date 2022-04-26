package main

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	repo "github.com/alexzh7/sample-service/internal/dvdstore/repository/postgres"
	"github.com/alexzh7/sample-service/internal/models"
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

	pgRepo := repo.NewPgRepo(db)

	prods := []*models.Product{
		{Id: 100, Quantity: 1},
		{Id: 5, Quantity: 3},
		{Id: 46, Quantity: 2},
	}
	ord, err := pgRepo.AddOrder(16, prods)

	if err == nil {
		fmt.Println(ord)
		for _, v := range ord.Products {
			fmt.Println(v)
		}
	}

}
