package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	db "github.com/HouzuoGuo/tiedot/db"
)

// Resources
// http://docopt.org/ - [docopt] Command-line interface description language
// https://medium.freecodecamp.org/writing-command-line-applications-in-go-2bc8c0ace79d How to write fast, fun command-line applications with Golang
// https://github.com/holman/boom/wiki/Commands
// https://zachholman.com/boom/ - motherfucking text snippets on the command line.

// Features:
// flipper - Base Command
// copy Item to clipboard
// add list
// add item into list
// remove item from list
// remove list
// List an item in a List

func main() {
	flag.Parse()

	var (
		listName  = flag.Arg(0)
		itemName  = flag.Arg(1)
		itemValue = flag.Arg(2)
	)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// open Database, create if not exists
	myDBDir := usr.HomeDir + "/.flipper"
	// myDBDir := "/.flipper"

	// (Create if not exist) open a database
	dbSess, err := db.OpenDB(myDBDir)
	if err != nil {
		fmt.Println("error opening the database")
		panic(err)
	}
	fmt.Println("database created or opened")

	// Trigger Action
	switch len(flag.Args()) {
	case 0:
		fmt.Println("Usage: flipper listname [itemname] [itemvalue]")
		os.Exit(1)
	case 1:
		if err := dbSess.Scrub(listName); err != nil {
			fmt.Println("list", listName, "exists")
			os.Exit(0)
		} else {
			fmt.Println("list", listName, "does not exists")
			if err := dbSess.Create(listName); err != nil {
				panic(err)
			}
			fmt.Println("list", listName, "created")
			os.Exit(0)
		}
	case 2:
		fmt.Println("case 2")
		if err := dbSess.Use(listName); err != nil {
			fmt.Println("list", listName, "exists")
			// Read document
			list := dbSess.Use(listName)
			fmt.Println("using list", listName)

			var query interface{}
			json.Unmarshal([]byte(`[{"eq": %itemName, "in": ["Title"]}]`), &query)

			queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

			if err := db.EvalQuery(query, list, &queryResult); err != nil {
				panic(err)
			}

			for id := range queryResult {
				// To get query result document, simply read it
				readBack, err := list.Read(id)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Query returned document %v\n", readBack)
				fmt.Println("copy", itemName, "from", listName, "into clipboard")
			}

			os.Exit(0)
		} else {
			fmt.Println("list", listName, "does not exists")
			if err := dbSess.Create(listName); err != nil {
				panic(err)
			}
			os.Exit(0)
		}

	case 3:
		fmt.Println("case 3")
		if err := dbSess.Use(listName); err != nil {
			fmt.Println("list", listName, "exists")

			list := dbSess.Use(listName)
			fmt.Println("using list", listName)

			var query interface{}
			json.Unmarshal([]byte(`[{"eq": %itemName, "in": ["Title"]}]`), &query)

			queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

			if err := db.EvalQuery(query, list, &queryResult); err != nil {
				panic(err)
			}
			fmt.Println("store queried for item")

			var itemExists bool

			for id := range queryResult {
				// To get query result document, simply read it
				readBack, err := list.Read(id)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Query returned document %v\n", readBack)
				itemExists = true
				fmt.Println("copy", itemName, "from", listName, "into clipboard")
			}

			fmt.Println(itemExists)

			if itemExists != true {

				fmt.Println("trying to insert data")
				docID, err := list.Insert(map[string]interface{}{
					"item":  itemName,
					"value": itemValue})
				if err != nil {
					panic(err)
				}
				fmt.Println("add item", itemName, "with value", itemValue, "into list", listName)

				// Read document
				readBack, err := list.Read(docID)
				if err != nil {
					panic(err)
				}
				fmt.Println("read item", itemName, "from list", listName)

				if err := list.Index([]string{"item"}); err != nil {
					panic(err)
				}
				fmt.Println("created index on item", itemName)

				fmt.Println(readBack)
			}
			os.Exit(0)
		} else {
			fmt.Println("list", listName, "does not exists")
			if err := dbSess.Create(listName); err != nil {
				panic(err)
			}
			fmt.Println("list", listName, "created")
			os.Exit(0)
		}
		fmt.Println("add value", itemValue, "of item", itemName, "into list", listName)

	}

}
