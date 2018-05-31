# flipper

snippet command line tool

## functions
- add list
- add item with value
- list all items of a list
- currently autooverwrite of values

## TODO
[ ] filepath for non windows systems  
[ ] multiline values  
[ ] copy to clib  
[ ] prompt for overwriting value  
[/] get lists  
[ ] do not break if file does not exist


## example file
```ini
[test]
test  = test2
test2 = test2

[test2]

[jokes]
chuck norris 1 = can coun't too much
chuck norris 2 = is afraid of ducks

[foo.bar]
test = value

[new List]
first Item = some value
```

## example create list


## example create item

## example add item to list
command: `flipper "new List" "first Item" "some value"`

<div style="font-family:consolas; color:#d7ba7d">item <span style="color:#0dbc79">first Item</span> with value <span style="color:#2172c8">some value</span> has been added to list <span style="color:#e5e511">new List</span></div>


## example get item


## example get items of list


## example get lists
feature of 2nd version