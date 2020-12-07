package xdb

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/nothollyhigh/kiss/log"
	kmysql "github.com/nothollyhigh/kiss/mysql"
)

var (
	dbMysql *kmysql.MysqlMgr
)

// DBConnInit init
func DBConnInit(conf kmysql.MgrConfig) {
	dbMysql = kmysql.NewMgr(conf)
	log.Info("DBConnInit Successfully !!!")
}

// RedisConnClose close
func RedisConnClose() {
	if dbMysql != nil {
		dbMysql.ForEach(func(tag string, idx int, m *kmysql.Mysql) {
			m.Close()
		})
	}
}

// OrmDBConn return gorm client
func OrmDBConn() *gorm.DB {
	if dbMysql == nil {
		return nil
	}
	return dbMysql.Get("master").OrmDB()
}

// DBConn return client
func DBConn(tag string) *sql.DB {
	if dbMysql == nil {
		return nil
	}
	return dbMysql.Get("master").DB()
}
