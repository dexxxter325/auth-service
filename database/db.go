package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *pgxpool.Pool

func ConnectToDb() (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file in db", err)
	}
	var (
		host     = os.Getenv("HOST")
		port     = os.Getenv("PORT")
		dbUser   = os.Getenv("USER")
		dbName   = os.Getenv("NAME")
		password = os.Getenv("PASSWORD")
		sslmode  = os.Getenv("SSLMODE")
	)

	//data := "host=localhost port=5432 user=postgres password=qwerty dbname=postgres sslmode=disable"
	data := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host,
		port,
		dbUser,
		dbName,
		password,
		sslmode,
	)
	DB, err = pgxpool.Connect(context.Background(), data)
	if err != nil {
		log.Fatalf("err in open pgx:%s", err)
	}
	return DB, nil
}
