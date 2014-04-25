package main

import (
    "github.com/eaigner/hood"
    "fmt"
    "os"
)

func main() {
    // Open a DB connection, use New() alternatively for unregistered dialects
    connection := os.Getenv("PQ_CONN")
    hd, err := hood.Open("postgres", connection)
    if err != nil {
        panic(err)
    }

    // Create a table
    type Fruit struct {
        Id    hood.Id
        Name  string `validate:"presence"`
        Color string
    }

    tx := hd.Begin()
    err = tx.CreateTable(&Fruit{})
    if err != nil {
        panic(err)
    }

    // Commit changes
    err = tx.Commit()
    if err != nil {
        panic(err)
    }

    fruits := []Fruit{
        Fruit{Name: "banana", Color: "yellow"},
        Fruit{Name: "apple", Color: "red"},
        Fruit{Name: "grapefruit", Color: "yellow"},
        Fruit{Name: "grape", Color: "green"},
        Fruit{Name: "pear", Color: "yellow"},
    }

    // Start a transaction
    tx = hd.Begin()

    ids, err := tx.SaveAll(&fruits)
    if err != nil {
        panic(err)
    }

    fmt.Println("inserted ids:", ids) // [1 2 3 4 5]

    // Commit changes
    err = tx.Commit()
    if err != nil {
        panic(err)
    }

    // Ids are automatically updated
    if fruits[0].Id != 1 || fruits[1].Id != 2 || fruits[2].Id != 3 {
        panic("id not set")
    }

    // If an id is already set, a call to save will result in an update
    fruits[0].Color = "green"

    ids, err = hd.SaveAll(&fruits)
    if err != nil {
        panic(err)
    }

    fmt.Println("updated ids:", ids) // [1 2 3 4 5]

    if fruits[0].Id != 1 || fruits[1].Id != 2 || fruits[2].Id != 3 {
        panic("id not set")
    }

    // Find
    //
    // The markers are db agnostic, so you can always use '?'
    // e.g. in Postgres they are replaced with $1, $2, ...
    var results []Fruit
    err = hd.Where("color", "=", "green").OrderBy("name").Limit(1).Find(&results)
    if err != nil {
        panic(err)
    }

    fmt.Println("results:", results) // [{1 banana green}]

    // Delete
    ids, err = hd.DeleteAll(&results)
    if err != nil {
        panic(err)
    }

    fmt.Println("deleted ids:", ids) // [1]

    results = nil
    err = hd.Find(&results)
    if err != nil {
        panic(err)
    }

    fmt.Println("results:", results) // [{2 apple red} {3 grapefruit yellow} {4 grape green} {5 pear yellow}]

    // Drop
    hd.DropTable(&Fruit{})
}
