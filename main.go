package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/HouzuoGuo/tiedot/db"
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

type myDB struct {
	db.DB
}

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
	os.RemoveAll(myDBDir)
	defer os.RemoveAll(myDBDir)

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		panic(err)
	}

	// Trigger Action
	switch len(flag.Args()) {
	case 0:
		fmt.Println("Usage: flipper listname [itemname] [itemvalue]")
		os.Exit(1)
	case 1:
		if listExists(listName) == true {
			fmt.Println("list ", listName, "exists")
			showAllItemsForList(listName)
		} else {
			fmt.Println("list ", listName, "does not exists")
			createList(listName)
		}
	case 2:
		copyItemToClipboard(listName, itemName)
	case 3:
		addItemToList(listName, itemName, itemValue)

	}

}

func copyItemToClipboard(listName string, itemName string) {

}

// func (db *badger.DB) writeToDB() {

// }

func readFromDB(listName string) {
	/*
		err := db.View(func(txn *badger.Txn) error {
			// Your code here…
			return nil
		})
	*/
}

func (myDB *myDB) listExists(listName string) bool {
	// Scrub (repair and compact) "Feeds"
	if err := myDB.Scrub(listName); err != nil {
		return false
	}
	return true
}

// flipper <list>
func (myDB *myDB) createList(listName string) {
	if err := myDB.Create(listName); err != nil {
		panic(err)
	}
}

// flipper <list> <name> <value>
func (myDB *myDB) addItemToList(listName string, itemName string, itemValue string) {
	feeds := myDB.Use(listName)
}

// flipper <list> remove <name> OR
// flipper remote <list> <name> OR
// flipper <list> <name> remove
func removeItemFromList() {

}

func (myDB *myDB) deleteList(listName string) {
	if err := myDB.Drop(listName); err != nil {
		panic(err)
	}
}

func deleteItem() {

}

func showAllItemsForList(listName string) {

}
