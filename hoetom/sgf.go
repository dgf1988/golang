package hoetom

import (
	"database/sql"
	"strings"
	"time"
	"fmt"
)

type Step struct {
	P int
	X int
	Y int
	C string
}

//常量，坐标索引
const stepch = "abcdefghijklmnopqrstuvwxyz"

//输出到Sgf
func (this Step) ToSgf() string {
	sgfitems := make([]string, 0)
	if this.X >= 0 && this.Y >= 0 && this.X < 19 && this.Y < 19 {
		if this.P == 1 {
			sgfitems = append(sgfitems, "B["+string(stepch[this.X])+string(stepch[this.Y])+"]")
		} else if this.P == 2 {
			sgfitems = append(sgfitems, "W["+string(stepch[this.X])+string(stepch[this.Y])+"]")
		}
	}
	if len(this.C) > 0 {
		sgfitems = append(sgfitems, "C["+this.C+"]")
	}
	if len(sgfitems) == 0 {
		return ""
	}
	return ";" + strings.Join(sgfitems, "")
}

//输出字符串
func (this Step) String() string {
	return this.ToSgf()
}

type Sgf struct {
	Id     int64
	Sgfid  int64
	Time   time.Time
	Place  string
	Event  string
	Black  string
	White  string
	Rule   string
	Result string
	Steps  string
	Update time.Time
}

func (this Sgf) ToSgf() string {
	items_sgf := make([]string, 0)
	items_sgf = append(items_sgf, "(;")
	items_sgf = append(items_sgf, fmt.Sprintf("EV[%s]", this.Event))
	if !this.Time.IsZero() {
		items_sgf = append(items_sgf, fmt.Sprintf("DT[%s]", this.Time.Format("2006-01-02")))
	} else {
		items_sgf = append(items_sgf, fmt.Sprintf("DT[0000-00-00]"))
	}
	items_sgf = append(items_sgf, fmt.Sprintf("PC[%s]", this.Place))
	items_sgf = append(items_sgf, fmt.Sprintf("PB[%s]", this.Black))
	items_sgf = append(items_sgf, fmt.Sprintf("PW[%s]", this.White))
	items_sgf = append(items_sgf, fmt.Sprintf("KO[%s]", this.Rule))
	items_sgf = append(items_sgf, fmt.Sprintf("RE[%s]", this.Result))
	items_sgf = append(items_sgf, "\n")
	items_sgf = append(items_sgf, this.Steps)
	items_sgf = append(items_sgf, ")")
	return strings.Join(items_sgf, "")
}

//输出链接
func (this Sgf) GetUrl() string {
	return UrlSgf(this.Sgfid)
}

//从数据库提取棋谱
func SgfGet(id int64) (*Sgf, error) {
	row := db.QueryRow("select * from sgf where id = ?", id)
	var get Sgf
	err := row.Scan(&get.Id, &get.Sgfid, &get.Time, &get.Place, &get.Event, &get.Black, &get.White, &get.Rule, &get.Result, &get.Steps, &get.Update)
	if err != nil {
		return nil, err
	}
	return &get, nil
}

//从数据库查找棋谱
func SgfFind(sid int64) (*Sgf, error) {
	row := db.QueryRow("select * from sgf where sid=?", sid)
	var find Sgf
	err := row.Scan(&find.Id, &find.Sgfid, &find.Time, &find.Place, &find.Event, &find.Black, &find.White, &find.Rule, &find.Result, &find.Steps, &find.Update)
	if err != nil {
		return nil, err
	}
	return &find, nil
}

//列出棋谱
func SgfList(take int, skip int) ([]Sgf, error) {
	rows, err := db.Query("select * from sgf order by id limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]Sgf, 0)
	for rows.Next() {
		var one Sgf
		err := rows.Scan(&one.Id, &one.Sgfid, &one.Time, &one.Place, &one.Event, &one.Black, &one.White, &one.Rule, &one.Result, &one.Steps, &one.Update)
		if err != nil {
			return list, err
		}
		list = append(list, one)
	}
	return list, rows.Err()
}


//向数据库添加棋谱
func SgfAdd(sgf *Sgf) (int64, error) {
	addsql := "insert into sgf (sid, stime, splace, sevent, sblack, swhite, srule, sresult, ssteps) values(?,?,?,?,?,?,?,?,?)"
	res, err := db.Exec(addsql, sgf.Sgfid, sgf.Time, sgf.Place, sgf.Event, sgf.Black, sgf.White, sgf.Rule, sgf.Result, sgf.Steps)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

//保存棋谱到数据库，如果数据库不存在这个棋谱
func SgfSave(sgf *Sgf) (int64, error) {
	find, err := SgfFind(sgf.Sgfid)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return find.Id, nil
	} else {
		return SgfAdd(sgf)
	}
}
