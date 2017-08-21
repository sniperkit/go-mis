package systemParameter

import (
	"log"
	"testing"
)

func TestFindByKey(t *testing.T) {
	// This test is assume that system parameter table is already exist data with key is equal to vt-backadate
	key := "vt-backdate"
	systemParameter, err := FindByKey(key)
	if err != nil {
		t.Fatal("Unable to execute query")
	}
	if systemParameter.Key != key {
		t.Fatal("Expected key: ", key, " actual key: ", systemParameter.Key)
	}
	log.Println(systemParameter)
	t.Log("Should return System Parameter")
}
