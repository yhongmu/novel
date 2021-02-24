package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

const (
	username	= "root"
	password	= "1017565696"
	net 	= "tcp"
	ip			= "121.4.83.168"
	port		= "3306"
	dbName		= "novel"
	driverName	= "mysql"
	charset		= "utf8"
)

var DB *sql.DB

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=uft8"
	//注意：要想解析time.Time类型，必须要设置parseTime=True
	//path := strings.Join([]string{username, ""}, "")
	path := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, net, ip, port, dbName, charset)
	var err error
	DB, err = sql.Open(driverName, path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	DB.SetConnMaxLifetime(100 * time.Second)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("database connect error, error info: " + err.Error())
		os.Exit(0)
	} else {
		fmt.Println("database connect success")
	}
}