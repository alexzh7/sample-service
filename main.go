package main

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/alexzh7/sample-service/models"
	"github.com/alexzh7/sample-service/repository"
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

	pgRepo := repository.NewPgRepo(db)

	prods := []*models.Product{
		{Id: 100, Quantity: 10},
	}
	ord, err := pgRepo.AddOrder(13, prods)

	if err == nil {
		fmt.Println(ord)
		for _, v := range ord.Products {
			fmt.Println(v)
		}
	}

	// TODO: LastInsertId is not supported by this driver

}
