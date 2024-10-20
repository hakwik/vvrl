package main

import "sort"

func sortedMapKeys(m interface{}) (keyList []string) {
	switch m := m.(type) {
	case map[string]string:
		for k := range m {
			keyList = append(keyList, k)
		}
	case map[string]bool:
		for k := range m {
			keyList = append(keyList, k)
		}
	default:
		panic("unknown map type")
	}

	sort.Strings(keyList)
	return
}
