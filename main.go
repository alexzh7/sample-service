package main

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/alexzh7/sample-service/models"
	_ "github.com/lib/pq"
)

func main() {
	//Logger
	l := zap.NewExample().Sugar()
	defer l.Sync()

	//TODO: Read config from env/file
	connStr := "user=pguser password=pgpass dbname=dvdstore sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal(err)
	}

	cst := models.NewCustomerPgRepo(db)
	customers, err := cst.GetCustomers(10)
	if err != nil {
		l.Fatal(err)
	}

	fmt.Println(customers)
	for _, v := range customers {
		fmt.Println(v.FirstName)
	}
}
