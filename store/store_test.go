package store

import (
	"testing"
	"time"
)

func TestStoreInitializeStore(t *testing.T) {
	var items map[int32]Item = make(map[int32]Item, 0)
	var expectedStore Store = Store{Counter: 0, Durations: 0, Items: items}

	if got := InitializeStore(); got.Counter != expectedStore.Counter || len(got.Items) != len(expectedStore.Items) {
		t.Errorf("Failed! %s : \n%v \ndoes not match expected : \n%v\n", "TestInitializeStore", got, expectedStore)
	}
}

func TestStoreCreateItem(t *testing.T) {

	var testStrings = []struct {
		tname    string
		strTest  string
		expected string
	}{
		{
			"Test Store Create Item",
			"angryMonkey",
			"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
		},
	}
	for _, test := range testStrings {
		if got, _ := CreateItem(test.strTest); &got != nil && got.Value != test.expected {
			t.Errorf("Failed! %s : \n%s \ndoes not match expected : \n%s\n", test.tname, got.Value, test.expected)
		}
	}
}

func TestStoreGetStats(t *testing.T) {

	var itemStore Store = Store{}
	itemStore.Counter = 2
	itemStore.Durations = 4000000

	var expectedStats Stats = Stats{Total: 2, Average: 2000000}

	var TestStoreGetStats = []struct {
		tname     string
		storeTest Store
		expected  Stats
	}{
		{
			"Test Store GetStats",
			itemStore,
			expectedStats,
		},
	}
	for _, test := range TestStoreGetStats {
		if got, err := test.storeTest.GetStats(); err != nil || got.Total != test.expected.Total || got.Average != test.expected.Average {
			t.Errorf("Failed! %s : \n%v \ndoes not match expected : \n%v\n", test.tname, got, test.expected)
		}
	}
}

func TestStoreGetItemById(t *testing.T) {

	var order int32 = 2
	var publishTime = time.Now()
	var unexpectedItem Item = Item{Value: "xyz", Publish: &publishTime}
	var expectedItem Item = Item{Value: "encryptedword", Publish: &publishTime}
	var items map[int32]Item = make(map[int32]Item, 0)
	items[0] = unexpectedItem
	items[1] = unexpectedItem
	items[order] = expectedItem

	var itemStore Store = Store{Items: items}

	var TestStoreGetItemById = []struct {
		tname     string
		storeTest Store
		expected  Item
	}{
		{
			"Test Store GetItemById",
			itemStore,
			expectedItem,
		},
	}
	for _, test := range TestStoreGetItemById {
		if got, err := test.storeTest.GetItemById(order); err != nil || got.Value != test.expected.Value || got.Publish != test.expected.Publish {
			t.Errorf("Failed! %s : \n%v \ndoes not match expected : \n%v\n", test.tname, got, test.expected)
		}
	}
}

func TestStoreGetItemByIdButTooSoon(t *testing.T) {

	var key int32 = 2
	var publishTime = time.Now().Add(time.Hour)
	var unexpectedItem Item = Item{Value: "xyz", Publish: &publishTime}
	var expectedItem Item = Item{Value: "encryptedword", Publish: &publishTime}
	var items map[int32]Item = make(map[int32]Item, 0)
	items[0] = unexpectedItem
	items[1] = unexpectedItem
	items[key] = expectedItem

	var itemStore Store = Store{Items: items}

	var TestStoreGetItemById = []struct {
		tname     string
		storeTest Store
		expected  Item
	}{
		{
			"Test Store GetItemByIdButTooSoon",
			itemStore,
			Item{},
		},
	}
	for _, test := range TestStoreGetItemById {
		//Expect Error
		if got, err := test.storeTest.GetItemById(key); err == nil || got != test.expected {
			t.Errorf("Failed! %s : \n%v \ndoes not match expected : \n%v\n", test.tname, got, test.expected)
		}
	}
}

func TestStoreStoreItem(t *testing.T) {

	var items map[int32]Item = make(map[int32]Item, 0)
	var itemStore Store = Store{Counter: 0, Durations: 0, Items: items}
	var requestTime = time.Now()
	var publishTime = requestTime.Add(time.Second * 5)
	var newItem Item = Item{Value: "encryptedword", Requested: &requestTime}
	var expectedItem Item = Item{Value: "encryptedword", Requested: &requestTime, Publish: &publishTime}
	var expectedKey int32 = 1

	// Call the function to add the newItem to itemStore
	// We should get an int32 key back

	if got, err := itemStore.StoreItem(newItem); err != nil || got != expectedKey {
		t.Errorf("Failed! %s : \n%v \ndoes not match expected : \n%v\n", "TestStore StoreItem", got, expectedKey)
	}

	// use the key returned to get the item from the itemStore
	// the itemValues should match expectedItem

	if got, ok := itemStore.Items[expectedKey]; !ok || *got.Publish != *expectedItem.Publish {
		t.Errorf("Failed! %s : \nok? %v\n%v \ndoes not match expected : \n%v\n", "TestStore Retrieve Item by Key after StoreItem", ok, got, expectedItem)
	}
}
