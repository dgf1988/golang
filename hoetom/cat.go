package hoetom

import (
	"database/sql"
)

type Cat struct {
	Id   int64
	Name string
}



//取国家
func CatGet(id int64) (*Cat, error) {
	var cat Cat
	row := db.QueryRow("SELECT * FROM cat WHERE id = ?", id)
	err := row.Scan(&cat.Id, &cat.Name)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

//找国家
func CatFind(name string) (*Cat, error) {
	var cat Cat
	row := db.QueryRow("SELECT * FROM cat WHERE name = ?", name)
	err := row.Scan(&cat.Id, &cat.Name)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func CatList(take, skip int) ([]Cat, error) {
	rows, err := db.Query("SELECT * FROM cat ORDER BY id LIMIT ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listcat := make([]Cat, 0)
	for rows.Next() {
		var cat Cat
		err := rows.Scan(&cat.Id, &cat.Name)
		if err != nil {
			return listcat, err
		}
		listcat = append(listcat, cat)
	}
	return listcat, rows.Err()
}

//添加国家
func CatAdd(name string) (int64, error) {
	res, err := db.Exec("INSERT INTO cat VALUES (?,?)", nil, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func CatDel(id int64) (int64, error) {
	res, err := db.Exec("DELETE FROM cat WHERE id = ?", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//保存国家
func CatSave(name string) (int64, error) {
	find, err := CatFind(name)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return find.Id, nil
	} else {
		return CatAdd(name)
	}
}

func CatRemove(name string) (int64, error) {
	find, err := CatFind(name)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return CatDel(find.Id)
	}
	return 0, nil
}

func CatUpdate(oldname, newname string) (int64, error) {
	res, err := db.Exec("UPDATE cat SET name = ? WHERE name = ?", newname, oldname)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

