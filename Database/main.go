package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://FSAE_DB_User@localhost:5432/telemetrydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), "insert into telemetry(date) values (NOW())")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insertion failed: %v\n", err)
		os.Exit(1)
	}

	rows, err := dbpool.Query(context.Background(), "select date from telemetry")

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	for rows.Next() {
		var date pgtype.Date
		oErr := rows.Scan(&date)
		if oErr != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", oErr)
			os.Exit(1)
		}
		fmt.Println(date)
	}	
}