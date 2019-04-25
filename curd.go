package mysql_unit

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mysql"
)

func initDB(config Config) *sql.DB {

	path := config.UserName + ":" + config.Password + "@tcp(" + config.IP + ":" + config.PORT + ")/" + config.DBName + "?charset=utf8"

	conn, err := sql.Open("mysql", path)
	if err != nil {
		checkErr(err)
	}
	conn.SetConnMaxLifetime(100)
	conn.SetMaxIdleConns(10)
	if err := conn.Ping(); err != nil {
		checkErr(err)
	}
	fmt.Println("connnect success")

	return conn

}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type DB struct {
	con *sql.DB
}

func New(config Config) *DB {
	con := initDB(config)
	return &DB{
		con,
	}
}

type field struct {
	Name       string
	Tag        string
	Type       string
	Key        bool
	Addr       interface{}
	IntSave    sql.NullInt64
	StringSave sql.NullString
	FloatSave  sql.NullFloat64
	BoolSave   sql.NullBool
}

type FieldsMap struct {
	dataobj interface{}
	reftype reflect.Type
	fields  []field
	table   string
	db      *sql.DB
}

func (obj *FieldsMap) GetFields() []field {
	return obj.fields
}

func newFieldsMap(table string, dataobj interface{}) (*FieldsMap, error) {

	elem := reflect.ValueOf(dataobj).Elem()
	reftype := elem.Type()

	var fields []field
	for i, flen := 0, reftype.NumField(); i < flen; i++ {

		var field field
		field.Type = reftype.Field(i).Type.String()
		field.Name = reftype.Field(i).Name
		field.Tag = reftype.Field(i).Tag.Get("sql")
		field.Addr = elem.Field(i).Addr().Interface()

		if reftype.Field(i).Tag.Get("key") == "" {
			field.Key = false
		} else {
			field.Key = true
		}
		fields = append(fields, field)
	}

	return &FieldsMap{
		dataobj: dataobj,
		reftype: reftype,
		fields:  fields,
		table:   table,
	}, nil
}

func (c *DB) NewFieldsMap(table string, dataobj interface{}) (*FieldsMap, error) {
	nfm, _ := newFieldsMap(table, dataobj)
	nfm.db = c.con
	return nfm, nil
}

func (fds *FieldsMap) GetFieldValues() []interface{} {

	var values []interface{}
	for i, flen := 0, len(fds.fields); i < flen; i++ {
		values = append(values, fds.GetFieldValue(i))
	}

	return values
}

func (fds *FieldsMap) GetFieldValue(idx int) interface{} {

	switch fds.fields[idx].Type {
	case "int64":
		return *fds.fields[idx].Addr.(*int64)
	case "string":
		return *fds.fields[idx].Addr.(*string)
	case "float64":
		return *fds.fields[idx].Addr.(*float64)
	case "bool":
		return *fds.fields[idx].Addr.(*bool)
	default:
	}

	return nil
}

func (c *FieldsMap) SQLFieldsStr() string {

	var tagsStr string
	for i, flen := 0, len(c.fields); i < flen; i++ {
		if len(tagsStr) > 0 {
			tagsStr += ", "
		}
		newTag := strings.Replace(c.fields[i].Tag, ".", "`.`", -1)
		newTag = strings.Replace(newTag, " as ", "` as `", -1)
		tagsStr += "`"
		tagsStr += newTag
		tagsStr += "`"
	}
	if len(tagsStr) > 0 {
		tagsStr += " "
		tagsStr = " " + tagsStr
	}

	return tagsStr
}

func (obj *FieldsMap) GetFieldSaveAddrs() []interface{} {

	var addrs []interface{}
	for i, flen := 0, len(obj.fields); i < flen; i++ {
		addrs = append(addrs, obj.GetFieldSaveAddr(i))
	}

	return addrs
}

func (fds *FieldsMap) GetFieldSaveAddr(idx int) interface{} {

	switch fds.fields[idx].Type {
	case "int64":
		return &fds.fields[idx].IntSave
	case "string":
		return &fds.fields[idx].StringSave
	case "float64":
		return &fds.fields[idx].FloatSave
	case "bool":
		return &fds.fields[idx].BoolSave
	default:
	}

	return nil
}

func (fds *FieldsMap) MapBackToObject() interface{} {

	for i, flen := 0, len(fds.fields); i < flen; i++ {
		switch fds.fields[i].Type {
		case "int64":
			if fds.fields[i].IntSave.Valid {
				*fds.fields[i].Addr.(*int64) = fds.fields[i].IntSave.Int64
			}
			break
		case "string":
			if fds.fields[i].StringSave.Valid {
				*fds.fields[i].Addr.(*string) = fds.fields[i].StringSave.String
			}
			break
		case "float64":
			if fds.fields[i].FloatSave.Valid {
				*fds.fields[i].Addr.(*float64) = fds.fields[i].FloatSave.Float64
			}
			break
		case "bool":
			if fds.fields[i].BoolSave.Valid {
				*fds.fields[i].Addr.(*bool) = fds.fields[i].BoolSave.Bool
			}
			break
		default:
		}

	}

	return fds.dataobj
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func (c *DB) BrowseToSource(table string, sql string, dataobj interface{}) error {

	elem := reflect.Indirect(reflect.ValueOf(dataobj))
	reftype := elem.Type()

	elemobj := reflect.Indirect(reflect.New(reftype.Elem().Elem())).Addr()

	obj, _ := newFieldsMap(table, elemobj.Interface())
	con := c.con
	_sql := strings.Join([]string{"SELECT ", obj.SQLFieldsStr(), " FROM ", obj.table, sql}, "")

	rows, err := con.Query(_sql)
	if err != nil {
		return err
	}

	for rows.Next() {
		nobj := reflect.Indirect(reflect.New(reftype.Elem().Elem())).Addr()
		fieldsMap, err := newFieldsMap(obj.table, nobj.Interface())
		if err != nil {
			return err
		}

		err = rows.Scan(fieldsMap.GetFieldSaveAddrs()...)
		if err != nil {
			return err
		}

		fieldsMap.MapBackToObject()
		elem = reflect.Append(elem, nobj)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	deepCopy(dataobj, elem.Interface())

	return err
}
