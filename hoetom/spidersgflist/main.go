package main

import (
	"fmt"
	"hoetom"
	"log"
	"os"
	"sync"
	"flag"
)

const (
	NameErrFile = "errsgflist.log"
)

var errlog *log.Logger
var waiter sync.WaitGroup

func main() {
	//参数
	flagcoroutine := flag.Int("c", 3, "number of coroutine for spider")
	flaglength := flag.Int("l", 100, "number of pages for spider")
	flagbegin := flag.Int("b", 1, "begin of page for spider")
	flag.Parse()

	//日志
	errfile, err := os.Create(NameErrFile)
	if err != nil {
		log.Panic(err.Error())
	}
	defer errfile.Close()
	errlog = log.New(errfile, "", log.LstdFlags)

	//任务量
	ch_page := make(chan int64, *flaglength)
	for i := *flagbegin; i < *flaglength +*flagbegin; i++ {
		ch_page <- int64(i)
	}
	close(ch_page)

	//协程量
	for i := 0; i < *flagcoroutine; i++ {
		waiter.Add(1)
		go Work(i, ch_page)
	}

	//等待结束
	waiter.Wait()
}

func Work(name int, ch_page chan int64) {
	for {
		page, ok := <-ch_page
		if !ok {
			fmt.Println(name, "退出")
			waiter.Done()
			return
		}
		url := hoetom.UrlSgfListLast(page)
		text, code := hoetom.Get(url)
		if code == hoetom.ErrCode {
			errlog.Println(name, page, text)
			continue
		}
		if code != 200 {
			errlog.Println(name, page, hoetom.HtmlTitle(text))
			continue
		}
		allplayerid := hoetom.HtmlAllPlayerid(text)
		rs1, err := hoetom.PlayeridSaveMany(allplayerid)
		if err != nil {
			errlog.Println(name, page, err.Error())
		}
		allsgfid := hoetom.HtmlAllSgfid(text)
		rs2, err := hoetom.SgfidSaveMany(allsgfid)
		if err != nil {
			errlog.Println(name, page, err.Error())
		}
		fmt.Println(name, "\tpage\t", page, "\tplayerid\t", rs1, "\tsgfid\t", rs2)
	}
}
