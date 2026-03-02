package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

var (
	dataUrl = "https://www.vihtavuori.com/wp-content/themes/vihtavuori/sovellus_vihtavuori/relodata.json"

	verbose      bool
	powder       string
	bulletweight string
	bulletname   string
	manufacturer string
	dl           = false
	data         VvData
	timeout      time.Duration
)

func init() {
	flag.BoolVar(&verbose, "v", false, "use verbose logging")
	flag.StringVar(&powder, "p", "", "filter by powder type")
	flag.StringVar(&bulletweight, "w", "", "filter by bullet weight (in grains)")
	flag.StringVar(&bulletname, "b", "", "filter by bullet model name")
	flag.StringVar(&manufacturer, "m", "", "filter by bullet manufacturer")
	flag.BoolVar(&dl, "d", false, "download new data from Vihtavuori")
	flag.DurationVar(&timeout, "t", 15*time.Second, "download timeout (only used with -d flag)")
}

func main() {
	flag.Parse()
	var err error

	d, err := loadData()
	if err != nil {
		fmt.Println("failed to load data:", err)
		d, err = downloadAndSave()
		if err != nil {
			fmt.Println("failed to save data:", err)
			os.Exit(1)
		}
		fmt.Println("saved data to local storage")
	}

	if dl {
		fmt.Println("downloading new data")
		d, err = downloadAndSave()
		if err != nil {
			fmt.Println("failed to download and save data:", err)
		}
	}

	data = d

	arg := flag.Arg(0)

	switch arg {
	case "cartridges":
		listCartridges()
		os.Exit(0)
	case "manufacturers":
		listBullets(0)
		os.Exit(0)
	case "bullets":
		listBullets(1)
		os.Exit(0)
	}

	if verbose {
		fmt.Println("Reloading data version:", data.Info[0].Greate)
	}

	cartridgeId := data.cartridgeIdFromName(arg)

	reloads := data.ReloData.
		filterByCartridgeId(cartridgeId).
		filterByBulletWeight(bulletweight).
		filterByPowderType(powder).
		filterByBulletMfg(manufacturer).
		filterByBulletName(bulletname)

	printTable(reloads, verbose)
}

func downloadAndSave() (VvData, error) {
	var d VvData
	d, err := download()
	if err != nil {
		return d, err
	}

	err = saveData(d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func download() (VvData, error) {
	fmt.Println("downloading data from Vihtavuori...")

	var data VvData

	t0 := time.Now()
	cli := http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", dataUrl, nil)
	if err != nil {
		return data, err
	}
	res, err := cli.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	if verbose {
		fmt.Printf("downloaded %d bytes in %v\n", len(body), time.Since(t0))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func listCartridges() {
	names := make([]string, 0)

	for _, v := range data.CartridgeData {
		names = append(names, v.Cartridge)
	}

	println("CARTRIDGE")
	sort.Strings(names)

	for _, v := range names {
		fmt.Println(v)
	}
}

func listBullets(i int) {
	bmap := data.bulletMap()
	uniqueBullets := make(map[string]bool)

	for _, b := range bmap {
		if b[0] == "" {
			continue
		}
		uniqueBullets[b[i]] = false
	}

	if i == 0 {
		fmt.Println("BULLET MANUFACTURER")
	} else {
		fmt.Println("BULLET MODEL")
	}

	for _, bullet := range sortedMapKeys(uniqueBullets) {
		fmt.Println(bullet)
	}
}

func saveData(data VvData) error {
	var err error
	storageLocation, err := getStorageLocation()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(storageLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)

	return nil
}

func loadData() (VvData, error) {
	var data VvData
	storageLocation, err := getStorageLocation()
	if err != nil {
		return VvData{}, err
	}

	f, err := os.OpenFile(storageLocation, os.O_RDONLY, 0600)
	if err != nil {
		return VvData{}, err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return VvData{}, err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return VvData{}, err
	}

	return data, nil
}

func getStorageLocation() (string, error) {
	var storageLocation string

	switch runtime.GOOS {
	case "windows":
		localappdata := os.Getenv("APPDATA")
		if localappdata == "" {
			return "", fmt.Errorf("could not locate user's appdata directory")
		}
		storageLocation = filepath.Join(localappdata, "vvrl.json")
	default:
		home := os.Getenv("HOME")
		if home == "" {
			return "", fmt.Errorf("could not locate user's home directory")
		}
		cfgDir := filepath.Join(home, ".config", "vvrl")
		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			// does not exist
			err = os.MkdirAll(cfgDir, 0700)
			if err != nil {
				return "", err
			}
		}
		storageLocation = filepath.Join(cfgDir, "vvrl.json")
	}
	return storageLocation, nil
}
