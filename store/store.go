package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/simrie/go_demo.git/hasher"
)

/*
Item type has the requested order, resulting crypto value and duration to complete encryption
*/
type Item struct {
	//Order     int32      `json:"order"`
	Value     string     `json:"value"`
	Requested *time.Time `json:"requested"`
	Publish   *time.Time `json:"publish"`
}

/*
Store type is a map of Items stored by request order
*/
type Store struct {
	Items     map[int32]Item `json:"items"`
	Counter   int32          `json:"counter"`
	Durations int64          `json:"durations"`
}

/*
Stats type has the values to be returned when Store statistics are requested
*/
type Stats struct {
	Total   int32   `json:"total"`
	Average float64 `json:"average"`
}

func InitializeStore() (s *Store) {
	items := make(map[int32]Item, 0)
	var store Store = Store{Items: items, Counter: 0}
	return &store
}

/*
GetStats function returns the total and average duration of all Store Items
*/
func (s *Store) GetStats() (Stats, error) {
	// short-circuit to prevent division by 0
	if s.Counter == 0 || s.Durations == 0 {
		var stats Stats = Stats{Total: s.Counter, Average: 0}
		return stats, nil
	}
	average := float64(s.Durations) / float64(s.Counter)
	var stats Stats = Stats{Total: s.Counter, Average: average}
	return stats, nil
}

/*
CreateItem accepts a string and returns an Item
that contains a Value that is the hash of the string
*/
func CreateItem(str string) (Item, error) {
	var item Item = Item{}
	if str == "" {
		return item, errors.New("missing string to encode")
	}
	hasher := hasher.Encode
	value, err := hasher(str)
	if err != nil {
		return item, errors.New("missing encoded value")
	}
	requested := time.Now()
	item.Requested = &requested
	item.Value = value
	return item, nil
}

/*
GetItemById returns the item by Order key only if
the current time is after the item's Publish time
*/
func (s *Store) GetItemById(i int32) (Item, error) {
	item, ok := s.Items[i]
	if !ok {
		msg := fmt.Sprintf("Item %d not found.", i)
		return Item{}, errors.New(msg)
	}
	t := time.Now()
	if t.Before(*item.Publish) {
		msg := fmt.Sprintf("Item %d not ready.", i)
		return Item{}, errors.New(msg)
	}
	return item, nil
}

/*
StoreItem stores the item only if it has Value and Requested time
*/
func (s *Store) StoreItem(item Item) (int32, error) {

	if item.Requested == nil {
		return 0, errors.New("missing Requested")
	}
	if item.Value == "" {
		return 0, errors.New("missing value")
	}
	t := item.Requested
	// Do not allow access to item until 5 seconds from now
	publishTime := t.Add(time.Second * 5)
	item.Publish = &publishTime
	// Get a duration for the time to hash and update store
	duration := publishTime.Sub(*item.Requested)
	s.Counter = s.Counter + 1
	s.Durations = s.Durations + duration.Nanoseconds()
	s.Items[s.Counter] = item
	return s.Counter, nil
}
