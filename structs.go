package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Reloads []Relodata

type Relodata struct {
	CartridgeID        int    `json:"CartridgeID"`
	BulletID           int    `json:"BulletID"`
	BulletWeightGrams  string `json:"BulletWeightGrams"`
	BulletWeightGrains string `json:"BulletWeightGrains"`
	ColMM              string `json:"ColMM"`
	ColInc             string `json:"ColInc"`
	BarrelLength       any    `json:"BarrelLength"`
	BarrelLengthInc    string `json:"BarrelLengthInc"`
	PowderType         string `json:"PowderType"`
	StartloadGrams     string `json:"StartloadGrams"`
	StartloadGrains    string `json:"StartloadGrains"`
	StartloadMs        string `json:"StartloadMs"`
	StartLoadFps       string `json:"StartLoadFps"`
	MaxloadGrams       string `json:"MaxloadGrams"`
	MaxloadGrains      string `json:"MaxloadGrains"`
	MaxloadMs          string `json:"MaxloadMs"`
	MaxLoadFps         string `json:"MaxLoadFps"`
	LoadType           string `json:"LoadType"`
	State              int    `json:"State"`
}

type VvData struct {
	Description string `json:"Description"`
	Version     string `json:"Version"`
	Bulletdata  []struct {
		BulletID   int    `json:"BulletId"`
		Bulletname string `json:"Bulletname"`
		Bulletmfg  string `json:"Bulletmfg"`
	} `json:"Bulletdata"`
	CartridgeData []struct {
		CartridgeID    int    `json:"CartridgeID"`
		CartridgeOrder int    `json:"CartridgeOrder"`
		Cartridge      string `json:"Cartridge"`
		CartridgeGroup string `json:"CartridgeGroup"`
		Barrel         string `json:"Barrel"`
		Primers        string `json:"Primers"`
		Cases          string `json:"Cases"`
		HeaderData     string `json:"HeaderData"`
		FooterData     string `json:"FooterData"`
		ModifyDate     string `json:"ModifyDate"`
		UsedInNewRelos int    `json:"UsedInNewRelos"`
	} `json:"CartridgeData"`
	PowderData []struct {
		PowderType  string `json:"PowderType"`
		PowderOrder int    `json:"PowderOrder"`
	} `json:"PowderData"`
	ReloComments []struct {
		CartridgeID       int    `json:"CartridgeId"`
		Bulletid          int    `json:"bulletid"`
		BulletWeightGrams string `json:"BulletWeightGrams"`
		CommentIndex      int    `json:"CommentIndex"`
		CommentColumn     string `json:"CommentColumn"`
		Comment           string `json:"Comment"`
	} `json:"ReloComments"`
	GartridgeGroup []struct {
		ID             int    `json:"ID"`
		Name           string `json:"Name"`
		UsedInNewRelos int    `json:"UsedInNewRelos"`
	} `json:"GartridgeGroup"`
	Relodata Reloads `json:"Relodata"`
	Info     []struct {
		State             string `json:"State:"`
		Greate            string `json:"Greate"`
		UsedTime          string `json:"UsedTime"`
		RelosTotal        int    `json:"RelosTotal"`
		RelosHandled      int    `json:"RelosHandled"`
		RelosUnHandled    int    `json:"RelosUnHandled"`
		TotalNewRelolines int    `json:"Total New Relolines"`
	} `json:"Info"`
}

func (data VvData) bulletMap() map[int][]string {
	bullets := make(map[int][]string)

	for _, v := range data.Bulletdata {
		bullets[v.BulletID] = []string{v.Bulletmfg, v.Bulletname}
	}

	return bullets
}

func (data VvData) cartridgeMap() map[int]string {
	cartridges := make(map[int]string)

	for _, v := range data.CartridgeData {
		cartridges[v.CartridgeID] = v.Cartridge
	}
	return cartridges
}

func (data VvData) bulletNameFromId(id int) (string, string) {
	mfg := data.bulletMap()[id][0]
	name := data.bulletMap()[id][1]
	return mfg, name
}

func (data VvData) bulletMfgFromId(id int) string {
	return data.bulletMap()[id][0]
}

func (data VvData) cartridgeNameFromId(id int) string {
	return data.cartridgeMap()[id]
}

func (data VvData) cartridgeIdFromName(name string) int {
	for k, v := range data.cartridgeMap() {
		if v == name {
			return k
		}
	}
	return 0
}

func (data VvData) cartridgeIdFromString(name string) int {
	for _, v := range data.CartridgeData {
		if v.Cartridge == name {
			return v.CartridgeID
		}
	}
	return 0
}

func (reloads Reloads) filterByCartridgeId(cartridgeId int) Reloads {
	matches := make([]Relodata, 0)
	for _, v := range reloads {
		if v.CartridgeID == cartridgeId {
			matches = append(matches, v)
		}
	}
	return matches
}

func (reloads Reloads) filterByBulletWeight(grains string) Reloads {
	if grains == "" {
		return reloads
	}

	matches := make([]Relodata, 0)
	for _, v := range reloads {
		if v.BulletWeightGrains == grains {
			matches = append(matches, v)
		}
	}
	return matches
}

func (reloads Reloads) filterByPowderType(powderType string) Reloads {
	if powderType == "" {
		return reloads
	}

	matches := make([]Relodata, 0)
	for _, v := range reloads {
		if v.PowderType == powderType {
			matches = append(matches, v)
		}
	}
	return matches
}

func (reloads Reloads) filterByBulletMfg(mfg string) Reloads {
	if mfg == "" {
		return reloads
	}

	matches := make([]Relodata, 0)
	for _, v := range reloads {
		if data.bulletMfgFromId(v.BulletID) == mfg {
			matches = append(matches, v)
		}
	}
	return matches
}

func (reloads Reloads) filterByBulletName(bulletname string) Reloads {
	if bulletname == "" {
		return reloads
	}

	matches := make([]Relodata, 0)
	for _, v := range reloads {
		_, name := data.bulletNameFromId(v.BulletID)
		if name == bulletname {
			matches = append(matches, v)
		}
	}
	return matches
}

func printTable(reloads Reloads) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprint(w, "BULLET\tBULLET WEIGHT\tPOWDER\tCOL\tMIN\tMAX\tMIN M/S\tMAX M/S\n")
	for _, v := range reloads {
		mfg, name := data.bulletNameFromId(v.BulletID)
		str := fmt.Sprintf("%s %s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", mfg, name, v.BulletWeightGrains, v.PowderType, v.ColMM, v.StartloadGrains, v.MaxloadGrains, v.StartloadMs, v.MaxloadMs)
		fmt.Fprint(w, str)
	}
	w.Flush()
}
