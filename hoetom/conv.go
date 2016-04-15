package hoetom

import (
	"strconv"
	"github.com/dgf1988/mahonia"
)

var coding mahonia.Decoder

func init() {
	coding = mahonia.NewDecoder(ConfigGetDefault("htmlcoding", "gb18030"))
}

func Atoi(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		panic(err.Error())
	}
	return n
}

func Atoi64(num string) int64 {
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return n
}
