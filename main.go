package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	color "github.com/fatih/color"
	ini "gopkg.in/ini.v1"
)

// Resources
// http://docopt.org/ - [docopt] Command-line interface description language
// https://medium.freecodecamp.org/writing-command-line-applications-in-go-2bc8c0ace79d How to write fast, fun command-line applications with Golang
// https://github.com/holman/boom/wiki/Commands
// https://zachholman.com/boom/ - motherfucking text snippets on the command line.
// https://github.com/go-ini/ini Package ini provides INI file read and write functionality in Go. https://ini.unknwon.io
// http://ascii.co.uk/art/dolphin - for splash screen?

// Features:
// flipper - Base Command
// copy Item to clipboard
// add list
// add item into list
// remove item from list
// remove list
// List an item in a List

var myFolder = "\\.flipper"
var fileName = "flipper.ini"
var flipperSplash = "Flipper!"

var filePath, listName, itemName, itemValue string

// terminal colors
var cList = color.New(color.FgYellow).SprintFunc()
var cItem = color.New(color.FgGreen).SprintFunc()
var cValue = color.New(color.FgBlue).SprintFunc()

func setHomeDir() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	filePath = usr.HomeDir + "\\" + myFolder + "\\" + fileName
}

func main() {
	flag.Parse()

	listName = flag.Arg(0)
	itemName = flag.Arg(1)
	itemValue = flag.Arg(2)

	setHomeDir()

	//filename is the path to the json config file
	cfg, err := ini.Load(filePath)
	if err != nil {
		fmt.Println("No file available. creating file", fileName, "in Folder", myFolder)
		f, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Could not create the file", err)
		}
		if err = f.Close(); err != nil {
			fmt.Println("Could not close the file", err)
		}
		os.Exit(0)
	} else {
		// fmt.Println("file loaded") // debug
	}

	// Trigger Action
	switch len(flag.Args()) {

	case 0:
		// TODO: check if file empty, if not we could print an overview of lists(count of items)

		lists := cfg.SectionStrings()
		if len(lists) > 0 {

			for _, nameOfList := range lists {
				sec, _ := cfg.GetSection(nameOfList)
				items := sec.Keys()
				if nameOfList != "DEFAULT" { // filter out default list (used for values without a list)
					fmt.Fprintln(color.Output, cList(nameOfList), "("+cItem(len(items))+")")
				}
			}

		} else {

			fmt.Fprintln(color.Output, "Usage: flipper", cList("listname"), cItem("[itemname]"), cValue("[itemvalue]"))
		}

		os.Exit(0)

	case 1:

		sec, err := cfg.GetSection(listName)
		if err == nil {
			// fmt.Fprintln(color.Output, "list", cList(listName), "exists") // debug
			items := sec.KeysHash()
			if len(items) > 0 {
				// fmt.Fprintln(color.Output, "items in List", cList(listName)) // debug
				for item, value := range items {
					fmt.Fprintln(color.Output, cItem(item), "=", cValue(value))
				}
			} else {
				// look in every list for item
				lists := cfg.SectionStrings()
				for _, nameOfList := range lists {
					sec, _ := cfg.GetSection(nameOfList)
					items := sec.KeysHash()
					if len(items) > 0 {
						// fmt.Fprintln(color.Output, "items in List", cList(listName)) // debug
						for item, value := range items {
							if item == listName {
								fmt.Fprintln(color.Output, flipperSplash, "value", cValue(value), "of item", cItem(item), "copyied to clipboard")
								os.Exit(0)
							}
						}
					}
				}

				fmt.Fprintln(color.Output, flipperSplash, "list", cList(listName), "is empty")
			}
		} else {
			createList(cfg, 1)
		}
		os.Exit(0)

	case 2:
		_, err := cfg.GetSection(listName)
		if err == nil {
			// fmt.Fprintln(color.Output, "list", cList(listName), "exists") // debug
			_, err := cfg.Section(listName).GetKey(itemName)
			if err == nil {
				// fmt.Fprintln(color.Output, "item", cItem(itemName), "exists in List", cList(listName)) // debug
				val := cfg.Section(listName).Key(itemName).Value()
				fmt.Fprintln(color.Output, flipperSplash, "copied", cValue(val), "from", cItem(itemName), "to clipboard [TODO: implement that]")
			} else {
				fmt.Fprintln(color.Output, flipperSplash, cItem(itemName), "does not exist in", cList(listName))
			}
		} else {
			fmt.Fprintln(color.Output, flipperSplash, cList(listName), "does not exist. Cannot copy", cItem(itemName), "to clipboard")
		}
		os.Exit(0)

	case 3:

		_, err := cfg.GetSection(listName)
		if err == nil {
			// fmt.Fprintln(color.Output, "list", cList(listName), "exists") // debug
			createItemOrOverwrite(cfg)
		} else {
			createList(cfg, 3)
			createItemOrOverwrite(cfg)
		}
		os.Exit(0)

	}

}

func createList(cfg *ini.File, step int) {
	cfg.NewSection(listName)
	if step == 1 {
		fmt.Fprintln(color.Output, flipperSplash, "list", cList(listName), "created")
	}
	cfg.SaveTo(filePath)
}

func createItemOrOverwrite(cfg *ini.File) {
	key, err := cfg.Section(listName).GetKey(itemName)
	if err == nil {
		// fmt.Fprintln(color.Output, "item", cItem(itemName), "exists in List", cList(listName)) // debug
		overwriteValue(cfg, key)
	} else {
		addItemToList(cfg)
	}
}

func addItemToList(cfg *ini.File) {
	_, err := cfg.Section(listName).NewKey(itemName, itemValue)
	if err == nil {
		fmt.Fprintln(color.Output, flipperSplash, "added", cValue(itemValue), "as", cItem(itemName), "to", cList(listName))
	}

	writeFile(cfg)
}

// func addItem(cfg *ini.File) {
// 	_, err := cfg.Section("").NewKey(listName, itemName)
// 	if err == nil {
// 		fmt.Fprintln(color.Output, flipperSplash, "added item", cItem(listName), "with value", cValue(itemName))
// 	}
// 	writeFile(cfg)
// }

func writeFile(cfg *ini.File) {
	err := cfg.SaveTo(filePath)
	if err != nil {
		fmt.Println("cannot write file")
	}
}

func overwriteValue(cfg *ini.File, key *ini.Key) {
	key.SetValue(itemValue)
	fmt.Fprintln(color.Output, flipperSplash, cItem(itemName), "overwritten with", cValue(itemValue), "in", cList(listName))
	cfg.SaveTo(filePath)
}
