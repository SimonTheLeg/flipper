package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

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

var (
	myFolder      = "\\.flipper"
	fileName      = "flipper.ini"
	flipperSplash = "Flipper!"
)

var deleteFlag = flag.String("d", "**ListName**", "help message for delete flag")

// settings
type Setting struct {
	name  string
	value string
}

var ListOfSettings = [2]Setting{
	{name: "setting.prompt", value: "false"},
	{name: "setting.sort", value: "false"}}

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

	readAndWriteSettings(cfg)

	// processing flags
	// -d
	/*
			* for the moment flag -d is always processed, even when not present in the command
			* therefore we check in the first else statement if one of the args contains the delete Flag
		  *
			* multiple deleteFlags are overwriting each other
			* "-d test3 -d test2" will result in test2
	*/
	if *deleteFlag != "**ListName**" {
		fmt.Println("delete flag found", *deleteFlag) // debug
		if flag.NArg() == 0 {
			fmt.Println("number of arguments is 0. deleteFlag:", *deleteFlag) // debug
			_, err := cfg.GetSection(*deleteFlag)
			if err == nil {
				deleteList(cfg, *deleteFlag)
				os.Exit(0)
			} else {

				lists := cfg.SectionStrings()
				for _, nameOfList := range lists {
					sec, _ := cfg.GetSection(nameOfList)
					items := sec.KeyStrings()
					if len(items) > 0 {
						// fmt.Fprintln(color.Output, "items in List", cList(listName)) // debug
						for _, item := range items {
							if item == *deleteFlag {
								cfg.Section(nameOfList).DeleteKey(*deleteFlag)
								fmt.Fprintln(color.Output, flipperSplash, "deleted", cItem(*deleteFlag), "from", cList(nameOfList))
								os.Exit(0)
							}
						}
					}
				}

				fmt.Fprintln(color.Output, flipperSplash, "List", cList(*deleteFlag), "does not exist")
				os.Exit(0)
			}
		} else if flag.NArg() == 1 {
			// the first argument of the command has been taken as the value of the delete Flag
			// therefore the number of arguments is reduced by one
			deleteItem(cfg, listName, *deleteFlag)
			// fmt.Fprintln(color.Output, flipperSplash, "deleted", cItem(listName), "from List", cList(*deleteFlag), "/ TODO: implement")
			os.Exit(0)
		}
		// }

	} else {
		//
		for _, argument := range flag.Args() {
			if strings.Contains(argument, "-d") {
				if flag.NArg() == 2 {
					deleteList(cfg, listName)
					os.Exit(0)
				} else if flag.NArg() == 3 {
					deleteItem(cfg, listName, itemName)
					// fmt.Fprintln(color.Output, flipperSplash, "deleted", cItem(listName), "from List", cList(itemName), "/ TODO: implement")
					os.Exit(0)
				}

			}
		}

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

	default:
		fmt.Fprintln(color.Output, "Usage: flipper", cList("listname"), cItem("[itemname]"), cValue("[itemvalue]"))

	}
}

func readAndWriteSettings(cfg *ini.File) {

	fileChanged := false

	sec, _ := cfg.GetSection("")
	settingsFromFile := sec.KeyStrings()

	if len(settingsFromFile) > 0 {
		// fmt.Fprintln(color.Output, "items in List", cList(listName)) // debug

		for i, setting := range ListOfSettings {
			foundInArray := false

			for _, item := range settingsFromFile {
				if item == setting.name {
					ListOfSettings[i].value = cfg.Section("").Key(setting.name).Value()
					// fmt.Fprintln(color.Output, flipperSplash, "found & read setting", setting.name) // debug
					foundInArray = true
					break
				}
			}

			if !foundInArray {
				_, err := cfg.Section("").NewKey(setting.name, setting.value)
				// fmt.Fprintln(color.Output, flipperSplash, setting.name, "not found. writing to file") // debug
				if err != nil {
					fmt.Fprintln(color.Output, flipperSplash, "could not write", setting.name, err)
				} else {
					fileChanged = true
				}
			}

		}

	} else {
		fmt.Println("settings empty, writing new") // debug
		for _, setting := range ListOfSettings {
			fmt.Println(setting.name, "=", setting.value)
			_, err := cfg.Section("").NewKey(setting.name, setting.value)
			if err != nil {
				fmt.Fprintln(color.Output, flipperSplash, "could not write", setting.name, err)
			} else {
				fileChanged = true
			}

		}
	}

	if fileChanged {
		writeFile(cfg)
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
		fmt.Fprintln(color.Output, flipperSplash, "added", cItem(itemName), "with value", cValue(itemValue), "to List", cList(listName))
	}

	writeFile(cfg)
}

func overwriteValue(cfg *ini.File, key *ini.Key) {
	key.SetValue(itemValue)
	fmt.Fprintln(color.Output, flipperSplash, cItem(itemName), "overwritten with value", cValue(itemValue), "in List", cList(listName))
	cfg.SaveTo(filePath)
}

func deleteList(cfg *ini.File, list string) {

	lists := cfg.SectionStrings()
	for _, nameOfList := range lists {
		if nameOfList == list {

			for _, setting := range ListOfSettings {
				if setting.name == "setting.prompt" {
					if setting.value == "true" {
						fmt.Println("do you really want to delete", list+"?")
					}
				}
			}

			cfg.DeleteSection(list)
			fmt.Fprintln(color.Output, flipperSplash, "List", cList(list), "deleted")
			writeFile(cfg)
		} else {
			fmt.Fprintln(color.Output, flipperSplash, "List", cList(list), "does not exist")
			os.Exit(0)
		}
	}
}

func deleteItem(cfg *ini.File, list string, item string) {

	lists := cfg.SectionStrings()
	for _, nameOfList := range lists {
		if nameOfList == list {

			keyExists := cfg.Section(list).HasKey(item)
			if keyExists == true {

				for _, setting := range ListOfSettings {
					if setting.name == "setting.prompt" {
						if setting.value == "true" {
							fmt.Println("do you really want to delete", list+"?")
						}
					}
				}

				cfg.Section(list).DeleteKey(item)
				fmt.Fprintln(color.Output, flipperSplash, "delted", cItem(item), "from List", cList(list))
				writeFile(cfg)
			} else {
				fmt.Fprintln(color.Output, flipperSplash, "item", cItem(item), "not found in List", cList(list))
				os.Exit(0)
			}
		}

	}
}

func writeFile(cfg *ini.File) {
	err := cfg.SaveTo(filePath)
	if err != nil {
		fmt.Println("cannot write file")
	}
}
