package main

import (
  "database/sql"
  "context"
  _ "github.com/databricks/databricks-sql-go"
)

func main() {
  dsn := "token:<access-token>@<db-url>.azuredatabricks.net:443/sql/1.0/warehouses/<warehouse-id>"

  db, err := sql.Open("databricks", dsn)
  if err != nil {
    panic(err)
  }

  rows, err := db.QueryContext(context.Background(), "SELECT * FROM <db>.<schema>.<table> LIMIT 5")
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  for rows.Next() {

  }

  defer db.Close()

  if err := db.Ping(); err != nil {
    panic(err)
  }

  print("ok")
}
