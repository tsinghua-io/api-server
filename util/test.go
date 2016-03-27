package util

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

const (
	UserId = "2013011187"
)

var (
	Password = os.Getenv("thu_pass")
)

func AssertDeepEqual(t *testing.T, actual, expected interface{}) bool {
	if !reflect.DeepEqual(actual, expected) {
		actualJson, _ := json.Marshal(actual)
		expectedJson, _ := json.Marshal(expected)
		t.Errorf("Actual: %s, Expected: %s", actualJson, expectedJson)
		return false
	}
	return true
}
