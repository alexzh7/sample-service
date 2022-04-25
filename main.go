package main

import (
	"database/sql"

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
		{Id: 100, Quantity: 396},
		{Id: 10, Quantity: 5},
		{Id: 40, Quantity: 5},
		{Id: 740, Quantity: 5},
		{Id: 23, Quantity: 5},
		{Id: 346, Quantity: 5},
		{Id: 98, Quantity: 5},
		{Id: 4, Quantity: 279},
	}
	_, err = pgRepo.AddOrder(12, prods)

	// // if err == nil {
	// // 	for _, v := range order.Products {
	// 		fmt.Println(v)
	// 	}
	// }

	// ord, err := pgRepo.GetCustomerOrders(19887)
	// for _, v := range ord {
	// 	fmt.Println(v.Id)
	// 	for _, p := range v.Products {
	// 		fmt.Println(p)
	// 	}
	// }
	// fmt.Println(err)
}
