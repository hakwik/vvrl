# vvrl
Dumb tool to simply display reloading data from Vihtavuori.

## Installation
```CGO_ENABLED=0 go install -ldflags "-s -w" .```

## Usage
List names of cartridges:

```vvrl cartridges```

List names of bullet manufacturers

```vvrl manufacturers```

List bullet models

```vvrl bullets```

To list all loads for the given cartridge:

```vvrl "6,5 Creedmoor"```

You can filter the selection with the following flags:
```
-p    Powder type (example: "N150")
-w    Bullet weight, in grains (example: "120")
-b    Bullet model name (example: "Ballistic Tip")
-m    Bullet manufacturer (example: "Lapua")
```

### Sub-commands
```
cartridges    List all avaliable cartridges
manufacturers List all available bullet manufacturers (with VV typos and all)
bullets       List all available bullet models
```


## Reloading data
It comes with reloading data bundled. You can however chose to download new data
and use that for display. That is triggered by the -d flag:

```vvrl -d -m Barnes -w 127 "6,5 Creedmoor"```

*Note* that this data is currently not stored, so you need to use the -d flag
to always get the current VV data.

### Updating bundled reloading data
It is certainly possible to update the bundled reloading data so you don't have to
use the `-d`flag all the time.. Upon build, the rldata.json file is embedded in the executable.

You can replace that rldata.json
file with new contents you download 
[from Vihtavuori](https://www.vihtavuori.com/wp-content/themes/vihtavuori/sovellus_vihtavuori/relodata.json) and run `go build` or `go install` to build and/or install a new
executable.