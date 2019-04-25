package mysql_unit

import (
	"log"
	"reflect"
)

func (obj *FieldsMap) Browse(sql string) ([]interface{}, error) {
	conn := obj.db
	_sql := "SELECT " + obj.SQLFieldsStr() + " FROM " + obj.table + sql

	rows, err := conn.Query(_sql)
	if err != nil {
		log.Fatal(err)
	}
	var objs []interface{}
	for rows.Next() {
		nobj := reflect.New(obj.reftype).Interface()
		fieldsMap, err := newFieldsMap(obj.table, nobj)
		if err != nil {
			return objs, err
		}

		err = rows.Scan(fieldsMap.GetFieldSaveAddrs()...)

		if err != nil {
			return objs, err
		}

		fieldsMap.MapBackToObject()
		objs = append(objs, nobj)
	}

	if err := rows.Err(); err != nil {
		return objs, err
	}

	return objs, nil

}

func (obj *FieldsMap) View(id int) (interface{}, error) {
	conn := obj.db
	_sql := "SELECT " + obj.SQLFieldsStr() + " FROM " + obj.table + " where id = ? "

	row := conn.QueryRow(_sql, id)

	nobj := reflect.New(obj.reftype).Interface()
	fieldsMap, err := newFieldsMap(obj.table, nobj)
	if err != nil {
		return nobj, err
	}

	err = row.Scan(fieldsMap.GetFieldSaveAddrs()...)

	if err != nil {
		return nobj, err
	}

	fieldsMap.MapBackToObject()

	return nobj, nil

}

func (obj *FieldsMap) Insert() (int64, error) {
	conn := obj.db
	var vs string
	var tagsStr string
	var values []interface{}
	for i, flen := 0, len(obj.fields); i < flen; i++ {

		if !obj.fields[i].Key {
			if len(vs) > 0 {
				vs += ", "
			}
			vs += "?"

			if len(tagsStr) > 0 {
				tagsStr += ", "
			}
			tagsStr += "`"
			tagsStr += obj.fields[i].Tag
			tagsStr += "`"

			values = append(values, obj.GetFieldValue(i))
		}
	}

	if len(tagsStr) > 0 {
		tagsStr += " "
		tagsStr = " " + tagsStr
	}

	sqlstr := "INSERT INTO `" + obj.table + "` (" + tagsStr + ") " +
		"VALUES (" + vs + ")"
	tx, _ := conn.Begin()
	res, err := tx.Exec(sqlstr, values...)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return res.LastInsertId()
}

func (obj *FieldsMap) Update() (int64, error) {
	conn := obj.db
	var tagsStr string
	var whereSql string
	var keyVal int64 = 0
	var values []interface{}
	for i, flen := 0, len(obj.fields); i < flen; i++ {

		if obj.fields[i].Key {
			keyVal = obj.GetFieldValue(i).(int64)
			whereSql = " where `" + obj.fields[i].Tag + "` = ? "
		} else {
			if len(tagsStr) > 0 {
				tagsStr += ", "
			}
			tagsStr += "`"
			tagsStr += obj.fields[i].Tag
			tagsStr += "`"
			tagsStr += " = ?"

			values = append(values, obj.GetFieldValue(i))
		}
	}

	if keyVal == 0 {
		return 0, nil
	}

	values = append(values, keyVal)

	if len(tagsStr) > 0 {
		tagsStr += " "
		tagsStr = " " + tagsStr
	}

	sqlstr := "UPDATE `" + obj.table + "` SET " + tagsStr + whereSql

	tx, _ := conn.Begin()
	res, err := tx.Exec(sqlstr, values...)
	if err != nil {
		return 0, err
	}

	tx.Commit()

	return res.LastInsertId()
}

func (obj *FieldsMap) Remove() (int64, error) {
	conn := obj.db
	var whereSql string
	var keyVal int64 = 0
	for i, flen := 0, len(obj.fields); i < flen; i++ {

		if obj.fields[i].Key {
			keyVal = obj.GetFieldValue(i).(int64)
			whereSql = " where `" + obj.fields[i].Tag + "` = ? "
		}
	}

	if keyVal == 0 {
		return 0, nil
	}

	sqlstr := "DELETE FROM `" + obj.table + "`  " + whereSql

	tx, _ := conn.Begin()
	res, err := tx.Exec(sqlstr, keyVal)
	if err != nil {
		log.Fatal("Exec fail", err)
	}

	tx.Commit()

	return res.RowsAffected()
}

func (obj *FieldsMap) ViewToSource(id int) error {
	conn := obj.db
	_sql := "SELECT " + obj.SQLFieldsStr() + " FROM " + obj.table + " where id = ? "

	row := conn.QueryRow(_sql, id)
	err := row.Scan(obj.GetFieldSaveAddrs()...)

	if err != nil {
		return err
	}
	obj.MapBackToObject()

	return err
}
