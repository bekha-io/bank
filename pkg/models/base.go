package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

var Db *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	Db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return Db
}
