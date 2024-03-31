package store

import (
	"fmt"

	"github.com/shivamaravanthe/superSite-BE/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	var err error
	dsn := "root:niveus@123@tcp(127.0.0.1:3306)/superSiteDB?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		fmt.Printf("Errored connecting DB: %v\n", err.Error())
		return
	}

	DBl, err := DB.DB()
	if err != nil {
		fmt.Printf("Errored connecting DB: %v\n", err.Error())
		return
	}
	DBl.SetMaxIdleConns(10)
	DBl.SetMaxOpenConns(100)

	runMigrations(DB)
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&model.Users{})
	db.AutoMigrate(&model.PasswordStore{})
}
