package store

import (
	"errors"
	"fmt"
	"time"
)

/*
Item type has the requested order, resulting crypto value and duration to complete encryption
*/
type Item struct {
	Order    int32         `json:"order"`
	Value    string        `json:"value"`
	Duration time.Duration `json:"duration"`
	Publish  *time.Time    `json:"publish"`
}

/*
Store type is a map of Items stored by request order
*/
type Store struct {
	Items   map[int32]Item `json:"items"`
	Counter int32          `json:"counter"`
}

/*
Stats type has the values to be returned when Store statistics are requested
*/
type Stats struct {
	Total   int32   `json:"total"`
	Average float64 `json:"average"`
}

/*
GetStats function returns the total and average duration of all Store Items
*/
func (s *Store) GetStats() (Stats, error) {
	numerator := int64(0)
	for _, item := range s.Items {
		numerator = item.Duration.Nanoseconds()
	}
	average := float64(numerator) / float64(s.Counter)
	var stats Stats = Stats{Total: s.Counter, Average: average}
	return stats, nil
}

/*
GetItem returns the item by Order key only if
the current time is after the item's Publish time
*/
func (s *Store) GetItem(i int32) (Item, error) {
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
StoreItem stores the item only if it has Value, Duration and Publish time
*/
func (s *Store) StoreItem(item Item) (int32, error) {

	if item.Duration == 0 {
		return 0, errors.New("missing duration")
	}
	if item.Publish == nil {
		return 0, errors.New("missing Publish")
	}
	if item.Value == "" {
		return 0, errors.New("missing value")
	}
	s.Counter = s.Counter + 1
	item.Order = s.Counter
	s.Items[item.Order] = item
	return item.Order, nil
}
