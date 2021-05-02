package hasher

import (
	"crypto/sha512"
	"encoding/base64"
	"time"

	"github.com/simrie/go_demo.git/store"
)

/*
Encode function gets the next order number from the store,
encodes the string, and calls a goroutine to store the new item.
*/
func Encode(itemStore *store.Store, str string) (int32, error) {
	var item store.Item
	start := time.Now()
	h := sha512.Sum512([]byte(str))
	value := "{SHA512}" + base64.StdEncoding.EncodeToString(h[:])
	t := time.Now()
	duration := t.Sub(start)
	item.Value = value
	item.Duration = duration
	// Do not allow access to item until 5 minutes from now
	publishTime := t.Add(time.Minute * 5)
	item.Publish = &publishTime
	order, err := itemStore.StoreItem(item)
	if err != nil {
		return 0, err
	}
	return order, nil
}
