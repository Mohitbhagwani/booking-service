package db

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDB() (*sql.DB, error) {
    connectionString := "postgres://postgres:postgres@localhost:5433/user_management?sslmode=disable"
    conn, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }
    db = conn
    fmt.Println("Connected to the database")
    return db, nil
}
