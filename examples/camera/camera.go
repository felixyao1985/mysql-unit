package camera

import mu "mysql-unit"

var MYSQL_CONFIG = mu.SQL_Config{
	UserName: "root",
	Password: "root",
	IP:       "172.0.0.1",
	PORT:     "3306",
	DBName:   "test",
}

var DB = mu.New(MYSQL_CONFIG)

func New(camera interface{}) {
	println(camera)
}

func init() {
	println("Camera")
}
