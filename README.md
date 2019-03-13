###封装MYSQL 小工具

make_camera.go 可根据数据库 生成对应的camera文件。

使用超简单 包含了 增删改查，基础功能


```gotemplate
    每个camera 中的 包含了数据库的字段名 和字段类型，
    默认操作和返回整个表的数据，可自行增删 需要查询和返回的字段
    
    type Demo struct {
        ID          int64  `sql:"id" key:"PRIMARY"`
        INTEGRATION string `sql:"integration"`
    }
```