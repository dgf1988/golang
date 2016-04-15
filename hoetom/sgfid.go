package hoetom

import (
	"database/sql"
	"sync"
)

var l_sgfid sync.Mutex

type Sgfid struct {
	Id int64
	Sgfid int64
	Status int64
}


//取棋手ID
func SgfidGet(id int64) (*Sgfid, error) {
	var get Sgfid
	row := db.QueryRow("select * from sgfid where id = ?", id)
	err := row.Scan(&get.Id, &get.Sgfid, &get.Status)
	if err != nil {
		return nil, err
	}
	return &get, nil
}

//找棋手ID
func SgfidFind(sid int64) (*Sgfid, error) {
	var find Sgfid
	row := db.QueryRow("select * from sgfid where sid = ?", sid)
	err := row.Scan(&find.Id, &find.Sgfid, &find.Status)
	if err != nil {
		return nil, err
	}
	return &find, nil
}

//列出棋手ID
func SgfidList(take int, skip int) ([]Sgfid, error) {
	rows, err := db.Query("select * from sgfid order by id limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]Sgfid, 0)
	for rows.Next() {
		var one Sgfid
		err = rows.Scan(&one.Id, &one.Sgfid, &one.Status)
		if err != nil {
			return list, err
		}
		list = append(list, one)
	}
	return list, rows.Err()
}

func SgfidQuery(status int64, take, skip int) ([]Sgfid, error) {
	rows, err := db.Query("select * from sgfid where status = ? order by id limit ?,?", status, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]Sgfid, 0)
	for rows.Next() {
		var one Sgfid
		err = rows.Scan(&one.Id, &one.Sgfid, &one.Status)
		if err != nil {
			return list, err
		}
		list = append(list, one)
	}
	return list, rows.Err()
}

//添加棋手ID
func SgfidAdd(sid int64) (int64, error) {
	res, err := db.Exec("insert into sgfid values (?,?,?)", nil, sid, 0)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func SgfidDel(id int64) (int64, error) {
	res, err := db.Exec("delete from sgfid where id = ?", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//保存棋手ID
func SgfidSave(sid int64) (int64, error) {
	find, err := SgfidFind(sid)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return find.Id, nil
	} else {
		return SgfidAdd(sid)
	}
}

func SgfidRemove(sid int64) (int64, error) {
	find, err := SgfidFind(sid)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}else if find != nil {
		return SgfidDel(find.Id)
	} else {
		return 0, nil
	}
}

func SgfidSet(id , status int64) (int64, error) {
	res, err := db.Exec("update sgfid set status = ? where id = ?", status, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//批量添加棋手ID
func SgfidAddMany(list []int64) (int64, error) {
	l_sgfid.Lock()
	defer l_sgfid.Unlock()
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	stmt, err := tx.Prepare("insert into sgfid values (?,?,?)")
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	var affertrows int64
	for _, one := range list {
		res, err := stmt.Exec(nil, one, 0)
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

func SgfidSaveMany(list []int64) (int64, error) {
	l_sgfid.Lock()
	defer l_sgfid.Unlock()
	//开启事务
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	//插入
	stmt, err := tx.Prepare("insert into sgfid values (?,?,?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//查找
	findstmt, err := tx.Prepare("select id from sgfid where sid = ?")
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return 0, err
	}
	//影响行
	var rowsaffects int64
	//查找ID
	var findid int64
	for _, one := range list {
		//查找
		row := findstmt.QueryRow(one)
		err := row.Scan(&findid)
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
		res, err := stmt.Exec(nil, one, 0)
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
