# vvrl
Dumb tool to simply display reloading data from Vihtavuori.

## Installation
```CGO_ENABLED=0 go install -ldflags "-s -w" .```

## Usage
List names of cartridges:

```vvrl cartridges```

To list all loads for the given cartridge:

```vvrl "6,5 Creedmoor"```

You can filter the selection with the following flags:
```
-p    Powder type (example: "N150")
-w    Bullet weight, in grains (example: "120")
-b    Bullet model name (example: "Ballistic Tip")
-m    Bullet manufacturer (example: "Lapua")
```

### Reloading data
It comes with reloading data bundled. You can however chose to download new data
and use that for display. That is triggered by the -d flag:

```vvrl -d -m Barnes -w 127 "6,5 Creedmoor"```
