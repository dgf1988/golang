package hoetom

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Player struct {
	Id     int64
	Pid    int64
	Name   string
	Sex    int64
	Rank   int64
	Nat    int64
	Cat    int64
	Birth  time.Time
	Update time.Time
}

func (this Player) GetUrl() *url.URL {
	strurl := fmt.Sprintf("http://www.hoetom.com/playerinfor_2011.jsp?id=%d", this.Pid)
	purl, err := url.Parse(strurl)
	if err != nil {
		return nil
	}
	return purl
}

//获取棋手
func PlayerGet(id int) (*Player, error) {
	var one Player
	row := db.QueryRow("select * from player where id = ?", id)
	err := row.Scan(&one.Id, &one.Pid, &one.Name, &one.Sex, &one.Rank, &one.Nat, &one.Cat, &one.Birth, &one.Update)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

//查找棋手
func PlayerFind(pid int64) (*Player, error) {
	var one Player
	row := db.QueryRow("select * from player where pid = ?", pid)
	err := row.Scan(&one.Id, &one.Pid, &one.Name, &one.Sex, &one.Rank, &one.Nat, &one.Cat, &one.Birth, &one.Update)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

func PlayerFindBy(name string) (*Player, error) {
	var one Player
	row := db.QueryRow("select * from player where player.pname=?", name)
	err := row.Scan(&one.Id, &one.Pid, &one.Name, &one.Sex, &one.Rank, &one.Nat, &one.Cat, &one.Birth, &one.Update)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

//列出棋手
func PlayerList(take int, skip int) ([]Player, error) {
	rows, err := db.Query("select * from player order by id limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]Player, 0)
	for rows.Next() {
		var one Player
		err := rows.Scan(&one.Id, &one.Pid, &one.Name, &one.Sex, &one.Rank, &one.Nat, &one.Cat, &one.Birth, &one.Update)
		if err != nil {
			return nil, err
		}
		list = append(list, one)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}
	return list, nil
}

//添加棋手
func PlayerAdd(p *Player) (int64, error) {
	res, err := db.Exec("insert into player (pid, pname, psex, prank, pnat, pcat, pbirth) values (?,?,?,?,?,?,?)", p.Pid, p.Name, p.Sex, p.Rank, p.Nat, p.Cat, p.Birth)
	if err != nil {
		return -1, err
	}
	insertid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	p.Id = insertid
	return p.Id, nil
}

func PlayerDel(id int64) (int64, error) {
	res, err := db.Exec("delete from player where id = ?", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//保存棋手
func PlayerSave(p *Player) (int64, error) {
	find, err := PlayerFind(p.Pid)
	if err != nil && err != sql.ErrNoRows {
		//错误
		return -1, err
	} else if find != nil {
		//已经存在
		return find.Id, nil
	} else {
		//执行添加
		return PlayerAdd(p)
	}
}

func PlayerSaveBy(playerid, name, sex, rank, nat, cat, birth string) (int64, error) {
	pid, err := strconv.ParseInt(playerid, 10, 64)
	if err != nil {
		return 0, err
	}
	psex := SexFromString(sex).Value
	prank, err := RankSave(rank)
	if err != nil {
		return 0, err
	}
	pnat, err := CountrySave(nat)
	if err != nil {
		return 0, err
	}
	pcat, err := CatSave(cat)
	if err != nil {
		return 0, err
	}
	pbirth, err := time.Parse("2006-01-02", birth)
	if err != nil {
		pbirth = time.Time{}
	}
	var player Player
	player.Pid = pid
	player.Name = name
	player.Sex = psex
	player.Rank = prank
	player.Nat = pnat
	player.Cat = pcat
	player.Birth = pbirth
	return PlayerSave(&player)
}
