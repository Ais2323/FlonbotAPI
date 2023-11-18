package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type ignoreWord struct {
	word string
}

func Init() {
	database, err := gorm.Open(sqlite.Open("ReplyDatabase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = database
}

func HasIgnoreWord(text string) bool {
	var exists bool
	model := ignoreWord{}
	err := db.Table("ignore_word").
		Model(model).
		Select("count(*) > 0").
		Where("? LIKE \"%\" || word || \"%\"", text).
		Find(&exists).
		Error
	if err != nil {
		panic(err)
	}
	return exists
}
