package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

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

// var deleteFlag = flag.String("d", "**ListName**", "help message for delete flag")

// settings
type Setting struct {
	name  string
	value string
}

var ListOfSettings = [2]Setting{
	{name: "setting.prompt", value: "false"},
	{name: "setting.sort", value: "false"}}

var filePath, listName, itemName, itemValue string

var arrayOfFlags = [8]string{"d", "delete", "l", "list", "s", "search", "a", "all"}

type Flag struct {
	name string
	pos  int
}

type Item struct {
	list  string
	item  string
	value string
}

var foundFlag Flag

// terminal colors
var cFlag = color.New(color.FgMagenta).SprintFunc()
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
	flag.Parse() // because of this flag.Arg() will be 1 instead of 0 if an argument is provided

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
	checkForFlags()
	readArguments()

	if foundFlag.name != "" {
		processFlags(cfg)
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
			showCommandStructure()
		}

		os.Exit(0)

	case 1:

		sec, err := cfg.GetSection(listName)
		if err == nil {
			fmt.Fprintln(color.Output, "list", cList(listName), "exists") // debug
			items := sec.KeysHash()
			if len(items) > 0 {
				// fmt.Fprintln(color.Output, "items in List", cList(listName)) // debug
				for item, value := range items {
					fmt.Fprintln(color.Output, cItem(item), "=", cValue(value))
				}
			}
		} else {
			result, err := lookForItem(cfg, itemName)
			if err == nil {
				copyToClipboard(cfg, result.item)
			} else {
				createList(cfg, listName)
			}
		}
		os.Exit(0)

	case 2:

		_, err := cfg.GetSection(listName)
		if err == nil {
			// fmt.Fprintln(color.Output, "list", cList(listName), "exists") // debug
			_, err := cfg.Section(listName).GetKey(itemName)
			if err == nil {
				// fmt.Fprintln(color.Output, "item", cItem(itemName), "exists in List", cList(listName)) // debug
				copyToClipboard(cfg, itemName)
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
			createList(cfg, listName)
			createItemOrOverwrite(cfg)
		}
		os.Exit(0)

	default:
		showCommandStructure()

	}
}

func showCommandStructure() {
	fmt.Fprintln(color.Output, "Usage: flipper", cFlag("[flag(s)]"), cList("listname"), cItem("[itemname]"), cValue("[itemvalue]"))
}

func copyToClipboard(cfg *ini.File, item string) {
	// look in every list for item
	result, err := lookForItem(cfg, item)
	if err == nil {
		// fmt.Fprintln(color.Output, flipperSplash, "value", cValue(value), "of item", cItem(item), "copyied to clipboard")
		fmt.Fprintln(color.Output, flipperSplash, "copied", cValue(result.value), "from", cItem(result.item), "to clipboard [TODO: implement that]")
	} else {
		fmt.Fprintln(color.Output, flipperSplash, "item", cItem(item), "not found")
	}
}

func lookForItem(cfg *ini.File, searchedItem string) (Item, error) {
	lists := cfg.SectionStrings()
	for _, nameOfList := range lists {
		sec, _ := cfg.GetSection(nameOfList)
		items := sec.KeysHash()
		if len(items) > 0 {
			// fmt.Fprintln(color.Output, "items in List", cList(searchedItem)) // debug
			for item, value := range items {
				if item == searchedItem {
					// fmt.Fprintln(color.Output, flipperSplash, "item", cItem(item), "found") // debug
					return Item{nameOfList, item, value}, nil
					// break
				}
			}
		}
	}
	err := errors.New("Not found")
	return Item{"", "", ""}, err
}

func readAndWriteSettings(cfg *ini.File) {
	// fmt.Println("running readAndWriteSettings") // debug

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

func checkForFlags() {
	// fmt.Println("running checkForFlags") // debug
	for i, item := range flag.Args() {
		if i == 0 || i == len(flag.Args())-1 {
			for _, flag := range arrayOfFlags {
				if item == flag {
					fmt.Println("Flag", item, "found")
					foundFlag = Flag{name: item, pos: i}
				}
			}
		}
	}
}

func readArguments() {
	fmt.Println("running readArguments(" + strconv.Itoa(len(flag.Args())) + ")") // debug
	if foundFlag.name != "" {
		listName = flag.Arg(1)
		itemName = flag.Arg(2)
		itemValue = flag.Arg(3)
	} else {
		listName = flag.Arg(0)
		itemName = flag.Arg(1)
		itemValue = flag.Arg(2)
	}
}

func processFlags(cfg *ini.File) {
	fmt.Println("running processFlags") // debug
	fmt.Println(foundFlag.name)
	if foundFlag.name == "d" || foundFlag.name == "delete" {
		processDeleteFlag(cfg)
	}
}

func processDeleteFlag(cfg *ini.File) {
	fmt.Println("running processDeleteFlag")
	// got one argument beside the flag, delete list (or item if list does not exist)
	if len(flag.Args()) == 2 {

		_, err := cfg.GetSection(listName)
		if err == nil {
			if deleteList(cfg) == nil {
				os.Exit(0)
			} else {
				var list = ""
				var item = listName
				deleteItem(cfg, list, item)
				os.Exit(0)
			}
		} else {
			// since the list was not found, check if we can delete an item with that name
			var item = listName
			result, err := lookForItem(cfg, item)
			if err == nil {
				deleteItem(cfg, result.list, result.item)
				// cfg.Section(result.list).DeleteKey(result.item)
				// fmt.Fprintln(color.Output, flipperSplash, "deleted", cItem(result.item), "from", cList(result.list))
				os.Exit(0)
			} else {
				fmt.Fprintln(color.Output, flipperSplash, "List", cList(listName), "does not exist")
				os.Exit(0)
			}

		}

		// got two arguments beside the flag, delete item
	} else if len(flag.Args()) == 3 {

	}
}

func createList(cfg *ini.File, list string) {
	cfg.NewSection(list)
	fmt.Fprintln(color.Output, flipperSplash, "list", cList(list), "created")
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

func deleteList(cfg *ini.File) error {
	fmt.Println("running deleteList") // debug

	listExists := false
	lists := cfg.SectionStrings()
	for _, nameOfList := range lists {
		fmt.Println("list", listName, "array list", nameOfList)
		if nameOfList == listName {

			listExists = true
			promptUser(listName)
			cfg.DeleteSection(listName)
			fmt.Fprintln(color.Output, flipperSplash, "List", cList(listName), "deleted")
			writeFile(cfg)
		}
	}
	if !listExists {
		fmt.Fprintln(color.Output, flipperSplash, "List", cList(listName), "does not exist")
		return fmt.Errorf("List %q not found", listName)
	}
	return nil
}

func promptUser(item string) {
	for _, setting := range ListOfSettings {
		if setting.name == "setting.prompt" {
			if setting.value == "true" {
				fmt.Println("do you really want to delete", item+"?")
			}
		}
	}
}

func deleteItem(cfg *ini.File, list string, item string) {
	fmt.Println("running deleteItem") // debug
	keyExists := false

	if list != "" {
		keyExists = cfg.Section(list).HasKey(item)
	} else {
		lists := cfg.SectionStrings()
		for _, nameOfList := range lists {
			keyExists := cfg.Section(nameOfList).HasKey(item)
			if keyExists {
				list = nameOfList
				break
			}
		}
	}

	if keyExists {
		promptUser(item)
		cfg.Section(list).DeleteKey(item)
		fmt.Fprintln(color.Output, flipperSplash, "delted", cItem(item), "from List", cList(list))
		writeFile(cfg)
	} else {
		if list != "" {
			fmt.Fprintln(color.Output, flipperSplash, "item", cItem(item), "not found in List", cList(list))
		} else {
			fmt.Fprintln(color.Output, flipperSplash, "item", cItem(item), "not found in any List")
		}
	}

}

func writeFile(cfg *ini.File) {
	err := cfg.SaveTo(filePath)
	if err != nil {
		fmt.Println("cannot write file")
	}
}
