package dao

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/chenminjian/spider/conf"
	"github.com/chenminjian/spider/utils/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var mutex sync.RWMutex
var disldb *sql.DB = nil
var logger = log.GetLogger("dao")
var stmtMap map[string]*sql.Stmt

//数据库连接
func Connect() {
	var dsn string

	db_host := conf.Conf.DBInfo.DBHost
	db_port := conf.Conf.DBInfo.DBPort
	db_user := conf.Conf.DBInfo.DBUser
	db_pass := conf.Conf.DBInfo.DBPass
	db_name := conf.Conf.DBInfo.DBName

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", db_user, db_pass, db_host, db_port, db_name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Errorf("open mysql database error.")
		panic(err)
	}
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(time.Second * 9)
	disldb = db

	stmtMap = make(map[string]*sql.Stmt)
}

func DB() *sql.DB {
	return disldb
}

func DBClose() error {
	return disldb.Close()
}

func Prepare(query string) (*sql.Stmt, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if stmt, ok := stmtMap[query]; ok {
		return stmt, nil
	}

	stmt, err := disldb.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "db prepare error")
	}
	stmtMap[query] = stmt

	return stmt, nil
}
