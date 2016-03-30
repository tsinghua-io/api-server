package util

import (
	"encoding/json"
	"net/http"
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

func NewFloat32(f float32) *float32 {
	return &f
}

func ExpectStatus(t *testing.T, actual, expected int) bool {
	if actual != expected {
		t.Errorf("Incorrect status code: expected %s, got %s", http.StatusText(actual), http.StatusText(expected))
		return false
	}
	return true
}

func ExpectDeepEqual(t *testing.T, actual, expected interface{}) bool {
	if !reflect.DeepEqual(actual, expected) {
		actualJson, _ := json.Marshal(actual)
		expectedJson, _ := json.Marshal(expected)
		t.Errorf("Actual: %s, Expected: %s", actualJson, expectedJson)
		return false
	}
	return true
}
