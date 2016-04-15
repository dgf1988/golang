package hoetom

import (
	"testing"
)

func aTest_Url(t *testing.T) {
	for i:=0; i < 28 ; i++ {
		t.Log(UrlPlayerList(i))
	}
}
