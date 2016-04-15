package hoetom

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"reflect"
	"errors"
)

type IScaner interface {
	Scan(dest ...interface{}) error
}


var db *sql.DB

func init() {
	DbInit()
}

func DbInit() {
	database, err := DbConnect()
	if err != nil {
		panic(err.Error())
	}
	db = database
}

func GetDb() *sql.DB {
	return db
}

func ScanRowToStruct(scaner IScaner, object interface{}) error {
	vp := reflect.ValueOf(object)
	if vp.Kind() == reflect.Ptr {
		vp = vp.Elem()
	} else {
		return errors.New("orm: the object must be a point of struct")
	}
	if vp.Kind() != reflect.Struct {
		return errors.New("orm: the object must be a point of struct")
	}
	scans := make([]interface{}, vp.NumField())
	for i := 0; i < vp.NumField(); i++ {
		scans[i] = vp.Field(i).Addr().Interface()
	}
	return scaner.Scan(scans...)
}

func DbConnect() (*sql.DB, error) {
	var (
		driver 	   string = ConfigGetDefault("dbdriver", "mysql")
		username   string = ConfigGetDefault("dbusername", "root")
		password   string = ConfigGetDefault("dbpassword", "guofeng001")
		hostname   string = ConfigGetDefault("dbhostname", "localhost")
		port       string = ConfigGetDefault("dbport", "3306")
		dbname     string = ConfigGetDefault("dbname", "weiqi_hoetom")
		charset    string = ConfigGetDefault("dbcharset", "utf8")
	)

	database, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v", username, password, hostname, port, dbname, charset, true))
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}
	return database, nil
}

func DbDesc(dbname string) ([][6]string, error) {
	rows, err := db.Query(fmt.Sprint("desc ", dbname))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	descs := make([][6]string, 0)
	for rows.Next() {
		scans := make([]interface{}, 6)
		for i := range scans {
			scans[i] = new(sql.NullString)
		}
		err := rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		var descrow [6]string
		for i := range scans {
			sqlnullvalue := scans[i].(*sql.NullString)
			if sqlnullvalue.Valid {
				descrow[i] = sqlnullvalue.String
			} else {
				descrow[i] = ""
			}
		}
		descs = append(descs, descrow)
	}
	return descs, rows.Err()
}

//统计数据量
func DbCount(tablename string) (int64, error) {
	row := db.QueryRow("select count(*) as num from " + tablename)
	var num int64
	err := row.Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func DbCountBy(tablename string, where string) (int64, error) {
	row := db.QueryRow("select count(*) as num from " + tablename + " where " + where)
	var num int64
	err := row.Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

//保存数据
func DbUpdate(tablename string, id int64, datas map[string]interface{}) (int64, error) {
	sqlitems := make([]string, 0)
	sqlitems = append(sqlitems, "update", tablename, "set")
	keys := make([]string, 0)
	args := make([]interface{}, 0)
	for k, v := range datas {
		keys = append(keys, k+" = ?")
		args = append(args, v)
	}
	sqlitems = append(sqlitems, strings.Join(keys, ","))
	sqlitems = append(sqlitems, "where id = ?")
	args = append(args, id)
	updatesql := strings.Join(sqlitems, " ")
	res, err := db.Exec(updatesql, args...)
	if err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

//
func DbClear(tablename string) (int64, error) {
	res, err := db.Exec(fmt.Sprint("TRUNCATE table ", tablename))
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
