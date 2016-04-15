package hoetom

import (
	"testing"
)

func Test_Config_Print(t *testing.T) {
	url := UrlSgf(135883)
	text, code := Get(url)
	if code == 200 {
		sgf, _ := HtmlFindSgf(text)
		t.Log(sgf)
	}
}
