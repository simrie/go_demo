package hasher

import (
	"testing"
)

func TestHasherEncode(t *testing.T) {

	var testStrings = []struct {
		tname    string
		strTest  string
		expected string
	}{
		{
			"Test Encode angryMonkey",
			"angryMonkey",
			"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
		},
	}
	for _, test := range testStrings {
		if got, _ := Encode(test.strTest); &got != nil && got != test.expected {
			t.Errorf("Failed! %s : \n%s \ndoes not match expected : \n%s\n", test.tname, got, test.expected)
		}
	}
}
