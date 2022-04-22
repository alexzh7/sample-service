package main

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

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

	ord, err := pgRepo.GetCustomerOrders(19887)
	for _, v := range ord {
		fmt.Println(v.Id)
		for _, p := range v.Products {
			fmt.Println(p)
		}
	}
	fmt.Println(err)
}
