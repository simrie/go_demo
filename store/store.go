package store

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/simrie/go_demo.git/hasher"
)

/*
Item type has the requested order, resulting crypto value and duration to complete encryption
*/
type Item struct {
	Order     int32      `json:"order"`
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
	realCount := int32(len(s.Items))
	if realCount == 0 || s.Durations == 0 {
		var stats Stats = Stats{Total: s.Counter, Average: 0}
		return stats, nil
	}
	average := float64(s.Durations) / float64(realCount)
	var stats Stats = Stats{Total: realCount, Average: average}
	return stats, nil
}

/*
CreateItem accepts a string and timestamp and creates an Item
that contains a Value that is the hash of the string
and then passes the Item to StoreItem as a goroutine.
*/
func CreateItem(str string, requested *time.Time, orderKey int32) (Item, error) {
	var item Item = Item{}
	if str == "" {
		return item, errors.New("missing string to encode")
	}
	if orderKey == 0 {
		return item, errors.New("missing order key")
	}
	hasher := hasher.Encode
	value, err := hasher(str)
	if err != nil {
		return item, errors.New("missing encoded value")
	}
	item.Requested = requested
	item.Value = value
	item.Order = orderKey
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
StoreItem stores the item only if it has Value, Requested time and order key
*/
func (s *Store) StoreItem(item Item) error {

	if item.Requested == nil {
		return errors.New("missing Requested")
	}
	if item.Value == "" {
		return errors.New("missing value")
	}
	if item.Order == 0 {
		return errors.New("missing order key")
	}
	t := item.Requested
	// Do not allow access to item until 5 seconds from now
	publishTime := t.Add(time.Second * 5)
	item.Publish = &publishTime
	// Get a duration for the time to hash and update store
	duration := publishTime.Sub(*item.Requested)
	//s.Counter = s.Counter + 1
	s.Durations = s.Durations + duration.Nanoseconds()
	s.Items[item.Order] = item
	return nil
}

/*
Worker goroutine receives the requested info
updates the Store.counter to get the new order key
posts the order key to the channel before
proceeding with the hashing and storing.
Errors are logged but not returned or posted to any channel.
*/
func (s *Store) Worker(plainText string, requested *time.Time, orderCh chan int32) {
	// update the store counter and claim the orderKey
	s.Counter = s.Counter + 1
	orderKey := int32(s.Counter)

	// post the new orderKey to the channel
	orderCh <- orderKey
	log.Printf(`Worker sent orderKey %d to order channel `, orderKey)

	// create the item that contains the hashed value
	item, err := CreateItem(plainText, requested, orderKey)
	if err != nil {
		log.Printf("\nWorker CreateItem %d Failed: %v: ", orderKey, err.Error())
		return
	}

	// store the item
	err = s.StoreItem(item)
	if err != nil {
		log.Printf("\nWorker StoreItem %d Failed: %v: ", orderKey, err.Error())
		return
	}
}
