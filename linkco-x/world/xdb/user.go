package xdb

import (
	"github.com/nothollyhigh/kiss/log"
)

func loadUser(tag, uid string) {
	db := DBConn(tag)
	username, basicData := "", ""
	if err := db.QueryRow(`select username,basic_data from tmpdb.tmptab where id=?`, uid).Scan(&username, &basicData); err != nil {
		log.Error("Find: (%s, %s), error: %v", tag, uid, err)
	}
}
