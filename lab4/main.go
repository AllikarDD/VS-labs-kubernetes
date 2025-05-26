package main
import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

const (
    host     = "postgres-service"
    port     = 5432
    user     = "dbuser"
    password = "dbpassword"
    dbname   = "mydatabase"
)

func main() {

    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    result, err := db.Exec("insert into Products (model, company, price) values ('iPhone X', $1, $2)",
        "Apple", 72000)
    if err != nil{
        panic(err)
    }
    fmt.Println(result.LastInsertId())  // не поддерживается
    fmt.Println(result.RowsAffected())  // количество добавленных строк
}