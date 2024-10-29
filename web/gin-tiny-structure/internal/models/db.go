package models

import (
	"fmt"
	"gdemo/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	config := config.GetConfig()

	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if config.MySQL.Debug {
		db = db.Debug()
	}
}

func init() {
	tables := []any{
		&User{},
		&CreditCard{},
	}

	if err := db.AutoMigrate(tables...); err != nil {
		panic(err)
	}
}
