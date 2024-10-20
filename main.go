package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

var (
	//go:embed rldata.json
	f string

	dataUrl = "https://www.vihtavuori.com/wp-content/themes/vihtavuori/sovellus_vihtavuori/relodata.json"

	verbose      bool
	powder       string
	bulletweight string
	bulletname   string
	data         VvData
	manufacturer string
	dl           = false
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

	if dl {
		f, err = download()
		if err != nil {
			fmt.Println("Something bad happened when downloading data:", err.Error())
			os.Exit(1)
		}
	}

	err = json.Unmarshal([]byte(f), &data)
	if err != nil {
		panic("could not unmarshal json: " + err.Error())
	}

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

	reloads := data.Relodata.filterByCartridgeId(cartridgeId).filterByBulletWeight(bulletweight).filterByPowderType(powder).filterByBulletMfg(manufacturer).filterByBulletName(bulletname)

	printTable(reloads)
}

func download() (string, error) {
	fmt.Println("downloading data from Vihtavuori...")

	t0 := time.Now()
	cli := http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", dataUrl, nil)
	if err != nil {
		return "", err
	}
	res, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if verbose {
		fmt.Printf("downloaded %d bytes in %v\n", len(body), time.Since(t0))
	}

	return string(body), nil
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
