package podcast

import (
	"fmt"
	"log"
	"sort"
)

// ItemList ..
type ItemList []*Item

// Len is part of sort.Interface.
func (items ItemList) Len() int {
	return len(items)
}

// Swap is part of sort.Interface.
func (items ItemList) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (items ItemList) Less(i, j int) bool {
	return items[i].Weight() > items[j].Weight()
}

// UnmarshalYAML ..
func (items *ItemList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal yaml items into map for easy reading
	var mItems = map[string]*Item{}
	unmarshal(&mItems)

	// populate `items`
	for key, item := range mItems {
		item.Key = key
		item.ExtractKeyInfo()

		// log.Printf("------------[%s] %+v", key, item.Title)
		*items = append(*items, item)
	}

	sort.Reverse(items)

	return nil
}

// Fix ..
func (items ItemList) Fix() {
	log.Printf("ItemList Fix()...")

	for _, item := range items {
		item.Fix()
	}
}

// Validate channel
func (items ItemList) Validate() error {
	log.Printf("ItemList Validate()...")

	// collect all guids and check if there is no duplicates
	var guids []string

	// Weight is season and episode weight (must be unique)
	var weights []int

	for _, item := range items {
		if err := item.Validate(); err != nil {
			return err
		}

		if inSlice(item.GUID.Text, guids) {
			return fmt.Errorf("GUID must be unique amongst all items. Found duplicate: `%s` (%s)", item.GUID.Text, item.Key)
		}
		if inSliceInt(item.Weight(), weights) {
			return fmt.Errorf("Weight (season, episode) must be unique amongst all items. Found duplicate: `%s` (%s)", item.Key, item.Title)
		}

		guids = append(guids, item.GUID.Text)
		weights = append(weights, item.Weight())

	}

	return nil
}
