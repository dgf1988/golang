package main

import (
	"flag"
	"fmt"
	"hoetom"
	"log"
	"os"
	"sync"
)

const (
	NameErrFile = "errsgf.log"
)

var errlog *log.Logger
var waiter sync.WaitGroup

func main() {
	flag_c := flag.Int("c", 3, "number of coroutines of spider")
	flag_s := flag.Int("s", 0, "status of sgfid")
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
	allid, err := hoetom.SgfidQuery(int64(*flag_s), *flag_l, 0)
	if err != nil {
		errlog.Println(err.Error())
		return
	}

	//
	ch_sgfid := make(chan hoetom.Sgfid, len(allid))
	for _, sgfid := range allid {
		ch_sgfid <- sgfid
	}
	close(ch_sgfid)

	//协程量
	for i := 0; i < *flag_c; i++ {
		waiter.Add(1)
		go Work(i, ch_sgfid)
	}

	//
	waiter.Wait()
}

func Work(name int, ch_id chan hoetom.Sgfid) {
	for {
		sgfid, ok := <-ch_id

		if !ok {
			fmt.Println(name, "退出")
			waiter.Done()
			return
		}

		url := hoetom.UrlSgf(sgfid.Sgfid)
		text, code := hoetom.Get(url)
		if code == hoetom.ErrCode {
			errlog.Println(name, sgfid.Sgfid, text)
			continue
		}

		if code != 200 {
			errlog.Println(name, sgfid.Sgfid, hoetom.HtmlTitle(text))
			continue
		}

		sgf, err := hoetom.HtmlFindSgf(text)
		if err != nil {
			errlog.Println(name, sgfid.Sgfid, err.Error())
			continue
		}

		insertid, err := hoetom.SgfSave(sgf)
		if err != nil {
			errlog.Println(name, sgfid.Sgfid, err.Error())
			continue
		}

		_, err = hoetom.SgfidSet(sgfid.Id, int64(code))
		if err != nil {
			errlog.Println(name, sgfid.Sgfid, err.Error())
			continue
		}

		fmt.Println(name, insertid, sgf.Place, sgf.Event, sgf.Black, sgf.White)
	}
}
