package camera

import (
	"strconv"
	"strings"
)

type ADVERSE_EVENT struct {
	ID          int64  `sql:"id" key:"PRIMARY"`
	INTEGRATION string `sql:"integration"`
}

var tableADVERSE_EVENT = "adverse_event"

func (c *ADVERSE_EVENT) Browse(sql string, row int, start int) ([]ADVERSE_EVENT, error) {

	_sql := strings.Join([]string{sql, " limt ", strconv.Itoa(start), " , ", strconv.Itoa(row)}, "")
	objs, _ := c.BrowseAll(_sql)
	return objs, nil
}

func (c *ADVERSE_EVENT) BrowseAll(sql string) ([]ADVERSE_EVENT, error) {

	fm, _ := DB.NewFieldsMap(tableADVERSE_EVENT, c)
	items, _ := fm.Browse(sql)

	var objs []ADVERSE_EVENT
	for i, olen := 0, len(items); i < olen; i++ {
		objs = append(objs, *items[i].(*ADVERSE_EVENT))
	}
	return objs, nil
}

func (c *ADVERSE_EVENT) View(id int) (ADVERSE_EVENT, error) {

	fm, _ := DB.NewFieldsMap(tableADVERSE_EVENT, c)
	items, _ := fm.View(id)
	return *items.(*ADVERSE_EVENT), nil
}

func (c *ADVERSE_EVENT) Insert() (int64, error) {
	fm, _ := DB.NewFieldsMap(tableADVERSE_EVENT, c)
	return fm.Insert()
}

func (c *ADVERSE_EVENT) Update() (int64, error) {
	fm, _ := DB.NewFieldsMap(tableADVERSE_EVENT, c)
	return fm.Update()
}

func (c *ADVERSE_EVENT) Remove() (int64, error) {
	fm, _ := DB.NewFieldsMap(tableADVERSE_EVENT, c)
	return fm.Remove()
}
