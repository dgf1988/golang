package main

import (
	"hoetom"
	"log"
	"os"
	"sync"
	"fmt"
	"flag"
)

const (
	NameErrFile = "errplayer.log"
)

var errlog *log.Logger
var waiter sync.WaitGroup
var l sync.Mutex


func main() {
	flag_c := flag.Int("c", 3, "number of coroutines of spider")
	flag_s := flag.Int("s", 0, "status of playerid")
	flag_l := flag.Int("l", 1000, "length of task")
	flag.Parse()

	//
	errfile, err := os.Create(NameErrFile)
	if err != nil {
		log.Panic(err.Error())
	}
	defer errfile.Close()
	errlog = log.New(errfile, "", log.LstdFlags)

	//任务量
	allplayerid, err := hoetom.PlayeridQuery(int64(*flag_s), *flag_l, 0)
	if err != nil {
		errlog.Println(err.Error())
		return
	}

	//
	ch_playerid := make(chan hoetom.Playerid, len(allplayerid))
	for _, playerid := range allplayerid {
		ch_playerid <- playerid
	}
	close(ch_playerid)

	//协程量
	for i:=0; i<*flag_c; i++ {
		waiter.Add(1)
		go Work(i, ch_playerid)
	}

	//
	waiter.Wait()
}

func Work(name int, ch_playerid chan hoetom.Playerid) {
	for {
		playerid, ok := <- ch_playerid
		if !ok {
			fmt.Println(name, "退出")
			waiter.Done()
			return
		}
		url := hoetom.UrlPlayer(playerid.Playerid)
		text, code := hoetom.Get(url)
		if code == hoetom.ErrCode {
			errlog.Println(name, playerid.Playerid, text)
			continue
		}
		if code != 200 {
			errlog.Println(name, playerid.Playerid, hoetom.HtmlTitle(text))
			continue
		}
		datas, err := hoetom.HtmlFindPlayer(text)
		if err != nil {
			errlog.Println(name, playerid.Playerid, err.Error())
			continue
		}
		l.Lock()
		insertid, err := hoetom.PlayerSaveBy(datas[0], datas[1], datas[2], datas[3], datas[4], datas[5], datas[6])
		l.Unlock()
		if err != nil {
			errlog.Println(name, playerid.Playerid, err.Error())
			continue
		}
		_, err = hoetom.PlayeridSet(playerid.Id, int64(code))
		if err != nil {
			errlog.Println(name, playerid.Playerid, err.Error())
			continue
		}
		fmt.Println(name, insertid, datas)
	}
}
