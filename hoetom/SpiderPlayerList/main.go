package main

import (
	"hoetom"
	"log"
	"os"
	"fmt"
)

const (
	NameErrFile = "errplayerlist.log"
)

func main() {
	errfile, err := os.Create(NameErrFile)
	if err != nil {
		log.Panic(err.Error())
	}
	defer errfile.Close()

	errlog := log.New(errfile, "", log.LstdFlags)
	for i := 0; i < 28; i++ {
		url_listplayer := hoetom.UrlPlayerList(i)
		text, code := hoetom.Get(url_listplayer)
		fmt.Println(i, code, url_listplayer)
		if code == hoetom.ErrCode {
			errlog.Println(i, code, url_listplayer, text)
			continue
		}
		if code != 200 {
			errlog.Println(i, code, url_listplayer, hoetom.HtmlTitle(text))
			continue
		}
		allid := hoetom.HtmlAllPlayerid(text)
		rs, err := hoetom.PlayeridSaveMany(allid)
		if err != nil {
			errlog.Println(err.Error())
			continue
		}
		fmt.Println("rowsaffects:", rs)
	}
}
