package gorm_dialects

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MySqlDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(uri), &gorm.Config{})
}
