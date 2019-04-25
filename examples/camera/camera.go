package camera

import mu "github.com/mysql-unit"

var Config = mu.Config{
	UserName: "root",
	Password: "root",
	IP:       "172.0.0.1",
	PORT:     "3306",
	DBName:   "test",
}

var DB = mu.New(Config)

func New(camera interface{}) {
	println(camera)
}

func init() {
	println("Camera")
}
