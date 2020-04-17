package qvs

import (
	"runtime/debug"
	"testing"

	"github.com/qiniu/api.v7/v7/auth"
)

var (
	testAccessKey = "Ves3WTXC8XnEHT0I_vacEQQz-9jrJZxNExcmarzQ"
	testSecretKey = "eNFrLXKG3R8TJ-DJA9YiMjLwuEfQnw8krrDuZzoy"
)

func skipTest() bool {
	return testAccessKey == "" || testSecretKey == ""
}

func getTestManager() *Manager {
	mac := auth.New(testAccessKey, testSecretKey)
	return NewManager(mac, nil)
}
func noError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Fatalf("should be nil, err = %s", err.Error())
	}
}

func shouldBeEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		debug.PrintStack()
		t.Fatalf("should be equal, expect = %#v, but got  = %#v", a, b)
	}
}
