package hoetom

import (
	"database/sql"
)

type Rank struct {
	Id   int64
	Name string
}

//取段位
func RankGet(id int64) (*Rank, error) {
	var rank Rank
	row := db.QueryRow("select * from rank where id = ?", id)
	err := row.Scan(&rank.Id, &rank.Name)
	if err != nil {
		return nil, err
	}
	return &rank, nil
}

//找段位
func RankFind(name string) (*Rank, error) {
	var rank Rank
	row := db.QueryRow("select * from rank where name = ?", name)
	err := row.Scan(&rank.Id, &rank.Name)
	if err != nil {
		return nil, err
	}
	return &rank, nil
}

func RankList(take, skip int)  ([]Rank, error) {
	rows, err := db.Query("select * from rank order by id limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list_rank := make([]Rank, 0)
	for rows.Next() {
		var rank Rank
		err := rows.Scan(&rank.Id, &rank.Name)
		if err != nil {
			return list_rank, nil
		}
		list_rank = append(list_rank, rank)
	}
	return list_rank, rows.Err()
}

//添加段位
func RankAdd(name string) (int64, error) {
	res, err := db.Exec("insert into rank values (?,?)", nil, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func RankDel(id int64) (int64, error) {
	res, err := db.Exec("delete from rank where id = ?",id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//保存段位
func RankSave(name string) (int64, error) {
	find, err := RankFind(name)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return find.Id, nil
	} else {
		return RankAdd(name)
	}
}

func RankRemove(name string) (int64, error) {
	find, err := RankFind(name)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return RankDel(find.Id)
	} else {
		return 0, nil
	}
}

func RankUpdate(oldname, newname string) (int64, error) {
	res, err := db.Exec("UPDATE rank SET name = ? WHERE name = ?", newname, oldname)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
