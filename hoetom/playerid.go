package hoetom

import (
	"database/sql"
	"sync"
)

var l_playerid sync.Mutex

type Playerid struct {
	Id       int64
	Playerid int64
	Status   int64
}

//取棋手ID
func PlayeridGet(id int64) (*Playerid, error) {
	var playerid Playerid
	row := db.QueryRow("select * from playerid where id = ?", id)
	err := row.Scan(&playerid.Id, &playerid.Playerid, &playerid.Status)
	if err != nil {
		return nil, err
	}
	return &playerid, nil
}

//找棋手ID
func PlayeridFind(playerid int64) (*Playerid, error) {
	var find Playerid
	row := db.QueryRow("select * from playerid where pid = ?", playerid)
	err := row.Scan(&find.Id, &find.Playerid, &find.Status)
	if err != nil {
		return nil, err
	}
	return &find, nil
}

//列出棋手ID
func PlayeridList(take int, skip int) ([]Playerid, error) {
	rows, err := db.Query("select * from playerid order by id limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listplayerid := make([]Playerid, 0)
	for rows.Next() {
		var playerid Playerid
		err = rows.Scan(&playerid.Id, &playerid.Playerid, &playerid.Status)
		if err != nil {
			return listplayerid, err
		}
		listplayerid = append(listplayerid, playerid)
	}
	return listplayerid, rows.Err()
}

func PlayeridQuery(status int64, take, skip int) ([]Playerid, error) {
	rows, err := db.Query("select * from playerid where status = ? order by id limit ?,?", status, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listplayerid := make([]Playerid, 0)
	for rows.Next() {
		var playerid Playerid
		err = rows.Scan(&playerid.Id, &playerid.Playerid, &playerid.Status)
		if err != nil {
			return listplayerid, err
		}
		listplayerid = append(listplayerid, playerid)
	}
	return listplayerid, rows.Err()
}

//添加棋手ID
func PlayeridAdd(playerid int64) (int64, error) {
	res, err := db.Exec("insert into playerid values (?,?,?)", nil, playerid, 0)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func PlayeridDel(id int64) (int64, error) {
	res, err := db.Exec("delete from playerid where id = ?", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//保存棋手ID
func PlayeridSave(pid int64) (int64, error) {
	playerid, err := PlayeridFind(pid)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if playerid != nil {
		return playerid.Id, nil
	} else {
		return PlayeridAdd(pid)
	}
}

func PlayeridRemove(pid int64) (int64, error) {
	find, err := PlayeridFind(pid)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}else if find != nil {
		return PlayeridDel(find.Id)
	} else {
		return 0, nil
	}
}

func PlayeridSet(id , status int64) (int64, error) {
	res, err := db.Exec("update playerid set status = ? where id = ?", status, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//批量添加棋手ID
func PlayeridAddMany(listplayerid []int64) (int64, error) {
	l_playerid.Lock()
	defer l_playerid.Unlock()
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	stmt, err := tx.Prepare("insert into playerid values (?,?,?)")
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	var affertrows int64
	for _, playerid := range listplayerid {
		res, err := stmt.Exec(nil, playerid, 0)
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return affertrows, err
		}
		as, err := res.RowsAffected()
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return affertrows, err
		}
		affertrows += as
	}
	stmt.Close()
	return affertrows, tx.Commit()
}

func PlayeridSaveMany(list_playerid []int64) (int64, error) {
	l_playerid.Lock()
	defer l_playerid.Unlock()
	//开启事务
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	//插入
	stmt, err := tx.Prepare("insert into playerid values (?,?,?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//查找
	findstmt, err := tx.Prepare("select id from playerid where pid = ?")
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return 0, err
	}
	//影响行
	var rowsaffects int64
	//查找ID
	var findplayerid int64
	for _, playerid := range list_playerid {
		//查找
		row := findstmt.QueryRow(playerid)
		err := row.Scan(&findplayerid)
		//找到
		if err == nil {
			continue
		}
		//错误
		if err != sql.ErrNoRows {
			findstmt.Close()
			stmt.Close()
			tx.Rollback()
			return rowsaffects, err
		}
		//不存在
		//执行插入
		res, err := stmt.Exec(nil, playerid, 0)
		if err != nil {
			//错误
			findstmt.Close()
			stmt.Close()
			tx.Rollback()
			return rowsaffects, err
		}
		//读取返回值
		rs, err := res.RowsAffected()
		if err != nil {
			findstmt.Close()
			stmt.Close()
			tx.Rollback()
			return rowsaffects, err
		}
		rowsaffects += rs
	}
	findstmt.Close()
	stmt.Close()
	//提交事务
	err = tx.Commit()
	if err != nil {
		//提交失败
		tx.Rollback()
		return 0, err
	}
	return rowsaffects, nil
}
