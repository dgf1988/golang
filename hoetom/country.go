package hoetom

import (
	"database/sql"
)

type Country struct {
	Id   int64
	Name string
}

//取国家
func CountryGet(id int64) (*Country, error) {
	var country Country
	row := db.QueryRow("SELECT * FROM country WHERE id = ?", id)
	err := row.Scan(&country.Id, &country.Name)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

//找国家
func CountryFind(name string) (*Country, error) {
	var country Country
	row := db.QueryRow("SELECT * FROM country WHERE name = ?", name)
	err := row.Scan(&country.Id, &country.Name)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func CountryList(take, skip int) ([]Country, error) {
	rows, err := db.Query("SELECT * FROM country ORDER BY id LIMIT ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list_country := make([]Country, 0)
	for rows.Next() {
		var country Country
		err := rows.Scan(&country.Id, &country.Name)
		if err != nil {
			return list_country, err
		}
		list_country = append(list_country, country)
	}
	return list_country, rows.Err()
}

//添加国家
func CountryAdd(name string) (int64, error) {
	res, err := db.Exec("INSERT INTO country VALUES (?,?)", nil, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func CountryDel(id int64) (int64, error) {
	res, err := db.Exec("DELETE FROM country WHERE id = ?", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

//保存国家
func CountrySave(name string) (int64, error) {
	find, err := CountryFind(name)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return find.Id, nil
	} else {
		return CountryAdd(name)
	}
}

func CountryRemove(name string) (int64, error) {
	find, err := CountryFind(name)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	} else if find != nil {
		return CountryDel(find.Id)
	}
	return 0, nil
}

func CountryUpdate(oldname, newname string) (int64, error) {
	res, err := db.Exec("UPDATE country SET name = ? WHERE name = ?", newname, oldname)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
