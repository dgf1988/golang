package main

import (
	"flag"
	"fmt"
	"hoetom"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Id     int64
	Pid    int64
	Name   string
	Sex    string
	Rank   string
	Nat    string
	Birth  time.Time
	Update time.Time
}

func playerlist(w http.ResponseWriter, r *http.Request) {

}

func playerbyid(w http.ResponseWriter, r *http.Request) {

}

func playerhandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	path = strings.Trim(r.URL.Path, "/")
	paths := strings.Split(path, "/")
	log.Println(paths)
	if len(paths) == 2 {
		id, err := strconv.ParseInt(paths[1], 10, 64)
		log.Println(id, err)
		if err == nil {
			player, err := hoetom.PlayerFind(id)
			if err != nil {
				log.Fatalln(err.Error())
				http.NotFound(w, r)
				return
			}
			t, err := template.ParseFiles("hoetom/server/tpl/player.html")
			if err != nil {
				log.Fatalln(err.Error())
				http.NotFound(w, r)
				return
			}
			playermap := make(map[string]string)
			playermap["name"] = player.Name
			playermap["sex"] = hoetom.Sex{player.Sex}.CnString()
			playermap["rank"] = strconv.FormatInt(player.Rank, 10)
			playermap["birth"] = player.Birth.Format("2006-01-02")
			playermap["update"] = player.Update.Format("2006-01-02")
			t.Execute(w, playermap)
			return
		}
		player, err := hoetom.PlayerFindBy(paths[1])
		if err != nil {
			log.Fatalln(err.Error())
			http.NotFound(w, r)
			return
		}
		t, err := template.ParseFiles("hoetom/server/tpl/player.html")
		if err != nil {
			log.Fatalln(err.Error())
			http.NotFound(w, r)
			return
		}
		playermap := make(map[string]string)
		playermap["name"] = player.Name
		playermap["sex"] = hoetom.Sex{player.Sex}.CnString()
		playermap["rank"] = strconv.FormatInt(player.Rank, 10)
		playermap["birth"] = player.Birth.Format("2006-01-02")
		playermap["update"] = player.Update.Format("2006-01-02")
		t.Execute(w, playermap)
		return
	}

	t, err := template.ParseFiles("hoetom/server/tpl/playerindex.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	list_player, err := hoetom.PlayerList(100, 0)
	if err != nil {
		log.Println(err.Error())
	} else {
		list := make([][8]string, 0)
		for _, p := range list_player {
			var parr [8]string
			parr[0] = strconv.Itoa(int(p.Pid))
			parr[1] = p.Name
			parr[2] = hoetom.Sex{p.Sex}.CnString()
			parr[3] = strconv.Itoa(int(p.Rank))
			parr[4] = strconv.Itoa(int(p.Nat))
			parr[5] = strconv.Itoa(int(p.Cat))
			parr[6] = p.Birth.Format("2006-01-02")
			parr[7] = p.Update.Format("2006-01-02")
			list = append(list, parr)
		}
		t.Execute(w, list)
	}
}

func main() {
	p := flag.String("p", ":8080", "http port")
	flag.Parse()

	http.Handle("/player/", NewPlayerHandler("/player/"))
	http.Handle("/sgf/", NewSgfHandler("/sgf/"))
	http.Handle("/css/", http.FileServer(http.Dir("hoetom/server")))
	http.Handle("/js/", http.FileServer(http.Dir("hoetom/server")))
	http.Handle("/img/", http.FileServer(http.Dir("hoetom/server")))
	http.ListenAndServe(*p, nil)
}

func ParseInt(str_value string) (int64, bool) {
	num, err := strconv.ParseInt(str_value, 10, 64)
	if err != nil {
		return 0, false
	}
	return num, true
}

func ParseTime(t time.Time) string {
	if t.IsZero() {
		return "0000-00-00"
	}
	return t.Format("2006-01-02")
}

/*

type Mux struct {
	Handlers map[string]http.Handler
}

func NewMux() *Mux {
	return &Mux{Handlers: make(map[string]http.Handler)}
}

func (this Mux) Handle(h http.Handler) {

}

*/

type playerData struct {
	Id     int64
	Pid    int64
	Name   string
	Sex    int64
	Rank   string
	Nat    string
	Cat    string
	Birth  time.Time
	Update time.Time
}

type PlayerHandler struct {
	Pattern string
}

func NewPlayerHandler(pattern string) *PlayerHandler {
	return &PlayerHandler{Pattern: pattern}
}

func (this PlayerHandler) Page(w http.ResponseWriter, r *http.Request, page int64) {
	if page <= 0 {
		log.Println(r.Method, r.URL.String(), "page < 0")
		http.NotFound(w, r)
		return
	}
	total, err := hoetom.DbCountBy("player", "player.pbirth>'0000-00-00'")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	sql := "select player.id, player.pid, player.pname, player.psex, rank.name, country.name, cat.name, player.pbirth, player.pupdate from player, rank, country, cat where player.pbirth > '0000-00-00' and player.prank=rank.id and player.pnat=country.id and player.pcat=cat.id order by player.id limit ?,100"
	rows, err := hoetom.GetDb().Query(sql, (page-1)*100)
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	defer rows.Close()
	//

	HeaderInfo{
		Title:"围棋棋手列表",
		Description:"围棋棋手列表",
		Keywords:"围棋,棋手,列表",
	}.WriteTo(w, "hoetom/server/tpl/header.html")

	remainder := total % 100
	if remainder > 0 {
		total = total/100 + 1
	} else {
		total = total / 100
	}
	err = HtmlIndexPage(w, total, page)
	if err != nil {
		log.Println(err.Error())
	}

	list_player := make([][9]string, 0)
	for rows.Next() {
		var p playerData
		err := hoetom.ScanRowToStruct(rows, &p)
		if err != nil {
			log.Println(r.Method, r.URL.String(), err.Error())
			http.NotFound(w, r)
			return
		}
		var player [9]string
		player[0] = strconv.FormatInt(p.Id, 10)
		player[1] = strconv.FormatInt(p.Pid, 10)
		player[2] = p.Name
		player[3] = hoetom.Sex{p.Sex}.CnString()
		player[4] = p.Rank
		player[5] = p.Nat
		player[6] = p.Cat
		player[7] = ParseTime(p.Birth)
		player[8] = ParseTime(p.Update)
		list_player = append(list_player, player)
	}

	t_player, err := template.ParseFiles("hoetom/server/tpl/playerlist.html")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	t_player.Execute(w, list_player)

	FooterInfo{
		Author:"dgf1988",
		Email:"dgf1988@qq.com",
		CopyRight:"©2016 weiqi163.com",
		ICP:"闽ICP备14014166号-2",
	}.WriteTo(w, "hoetom/server/tpl/footer.html")
}

func (this PlayerHandler) GetOneByName(w http.ResponseWriter, r *http.Request, name string) {
	sql := "select player.id, player.pid, player.pname, player.psex, rank.name, country.name, cat.name, player.pbirth, player.pupdate from player, rank, country, cat where player.pname = ? and player.prank=rank.id and player.pnat=country.id and player.pcat=cat.id"
	row := hoetom.GetDb().QueryRow(sql, name)
	var p playerData
	err := hoetom.ScanRowToStruct(row, &p)
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}

	var player [9]string
	player[0] = strconv.FormatInt(p.Id, 10)
	player[1] = strconv.FormatInt(p.Pid, 10)
	player[2] = p.Name
	player[3] = hoetom.Sex{p.Sex}.CnString()
	player[4] = p.Rank
	player[5] = p.Nat
	player[6] = p.Cat
	player[7] = ParseTime(p.Birth)
	player[8] = ParseTime(p.Update)
	t_header, err := template.ParseFiles("hoetom/server/tpl/header.html")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	t_header.Execute(w, player[2])

	t_player, err := template.ParseFiles("hoetom/server/tpl/player.html")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	t_player.Execute(w, player)

	t_footer, err := template.ParseFiles("hoetom/server/tpl/footer.html")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	t_footer.Execute(w, nil)
}

func (this PlayerHandler) GetOneByPid(w http.ResponseWriter, r *http.Request, pid int64) {
	sql := "select player.id, player.pid, player.pname, player.psex, rank.name, country.name, cat.name, player.pbirth, player.pupdate from player, rank, country, cat where player.pid = ? and player.prank=rank.id and player.pnat=country.id and player.pcat=cat.id"
	row := hoetom.GetDb().QueryRow(sql, pid)
	var p playerData
	err := hoetom.ScanRowToStruct(row, &p)
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}

	var player [9]string
	player[0] = strconv.FormatInt(p.Id, 10)
	player[1] = strconv.FormatInt(p.Pid, 10)
	player[2] = p.Name
	player[3] = hoetom.Sex{p.Sex}.CnString()
	player[4] = p.Rank
	player[5] = p.Nat
	player[6] = p.Cat
	player[7] = ParseTime(p.Birth)
	player[8] = ParseTime(p.Update)

	HeaderInfo{
		Title:       fmt.Sprint(p.Name, " - 棋手"),
		Description: fmt.Sprint(player[5], player[6], p.Name),
		Keywords:    strings.Join([]string{p.Name, player[4], player[6], "围棋", "棋手"}, ","),
	}.WriteTo(w, "hoetom/server/tpl/header.html")

	t_player, err := template.ParseFiles("hoetom/server/tpl/player.html")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	t_player.Execute(w, player)

	FooterInfo{
		Author:"dgf1988",
		Email:"dgf1988@qq.com",
		CopyRight:"©2016 weiqi163.com",
		ICP:"闽ICP备14014166号-2",
	}.WriteTo(w, "hoetom/server/tpl/footer.html")
}


func NewIndexPage(total, now int64) *indexPage {
	ip := indexPage{Steps: make([]indexPageStep, total)}
	for i := int64(1); i <= total; i++ {
		ip.Steps[i-1].Number = i
		if i == now {
			ip.Steps[i-1].IsNow = true
		} else {
			ip.Steps[i-1].IsNow = false
		}
	}
	return &ip
}

func (this PlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())
	if r.Method == "POST" {
		http.NotFound(w, r)
		return
	}
	path := strings.TrimSpace(r.URL.Path)
	path = strings.Trim(r.URL.Path, "/")
	paths := strings.Split(path, "/")
	if len(paths) == 1 {
		this.Page(w, r, 1)
		return
	}
	if len(paths) == 2 {
		num, ok := ParseInt(paths[1])
		if ok {
			this.GetOneByPid(w, r, num)
			return
		}
		this.GetOneByName(w, r, paths[1])
		return
	}
	if len(paths) == 3 && paths[1] == "page" {
		page, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			log.Println(r.Method, r.URL.String(), err.Error())
			http.NotFound(w, r)
			return
		}
		this.Page(w, r, page)
		return
	}
	http.NotFound(w, r)
}

type SgfHandle struct {
	Pattern string
}

func NewSgfHandler(pattern string) *SgfHandle {
	return &SgfHandle{Pattern: pattern}
}

func (this SgfHandle) SgfById(w http.ResponseWriter, r *http.Request, id int64) {
	sgf, err := hoetom.SgfGet(id)
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}

	HeaderInfo{
		Title:sgf.Event,
	}.WriteTo(w, "hoetom/server/tpl/header.html")

	t_sgf, err := template.ParseFiles("hoetom/server/tpl/weiqi_sgf.html")
	if err != nil {
		log.Println(r.Method, r.URL.String(), err.Error())
		http.NotFound(w, r)
		return
	}
	t_sgf.Execute(w, sgf.ToSgf())

	FooterInfo{
		Author:"dgf1988",
		Email:"dgf1988@qq.com",
		CopyRight:"©2016 weiqi163.com",
		ICP:"闽ICP备14014166号-2",
	}.WriteTo(w, "hoetom/server/tpl/footer.html")
}

func (this SgfHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	path := strings.TrimSpace(r.URL.Path)
	path = strings.Trim(r.URL.Path, "/")
	paths := strings.Split(path, "/")
	if len(paths) == 2 {
		num, err := strconv.ParseInt(paths[1], 10, 64)
		if err != nil {
			log.Println(r.Method, r.URL.String(), err.Error())
			http.NotFound(w, r)
			return
		}
		this.SgfById(w, r, num)
		return
	}
	http.NotFound(w, r)
}
