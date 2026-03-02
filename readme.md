# vvrl
Dumb tool to simply display reloading data from Vihtavuori.
![Screenshot](/doc/screenshot.png)

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

Verbose output also provides efficiency and sensitivity metrics (see [this article](https://storage.googleapis.com/wzukusers/user-33607572/documents/627aa957424246239bb5fe1bbb1204bf/Load%20Development%20v2.pdf)):
```
-v    Verbose output
```

### Sub-commands
```
cartridges    List all avaliable cartridges
manufacturers List all available bullet manufacturers (with VV typos and all)
bullets       List all available bullet models
```


## Reloading data
If data has not been previously fetched, it will attempt to download data on first use.
You can force downloading of new data, replacing the current data. This is triggered by the -d flag:

```vvrl -d -m Barnes -w 127 "6,5 Creedmoor"```
