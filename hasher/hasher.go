package hasher

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
	"time"

	"github.com/simrie/go_demo.git/store"
)

/*
Encode function gets the next order number from the store,
encodes the string, and calls a goroutine to store the new item.
*/
func Encode(itemStore *store.Store, str string) (int32, error) {
	var item store.Item
	requested := time.Now()
	item.Requested = &requested

	h := sha512.Sum512([]byte(str))
	value := base64.StdEncoding.EncodeToString(h[:])
	log.Println(value)
	item.Value = value
	order, err := itemStore.StoreItem(item)
	if err != nil {
		log.Printf(`Hasher error %v`, err)
		return 0, err
	}
	return order, nil
}
