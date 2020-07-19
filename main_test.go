package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Test struct {
	gorm.Model
	Name string
}

func setup() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:1234@(172.17.0.2)/test?parseTime=true")
	if err != nil {
		return nil, err
	}
	db.DropTableIfExists(&Test{})
	db.AutoMigrate(&Test{})
	return db, nil
}

func populate(db *gorm.DB) {
	for i := 0; i < 100; i++ {
		db.Model(&Test{}).Create(&Test{
			Name:  fmt.Sprintf("Entity-%d", i),
		})
	}
}



