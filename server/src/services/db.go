package services

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type DB struct {
	Db *gorm.DB
}

func (db DB) Insert(record interface{}) (err error) {
	result := db.Db.Create(record)

	if result.Error != nil {
		err = errors.New("insert failed")
	}

	return
}

func (db DB) Find(model interface{}) (err error) {
	fmt.Println(model)
	db.Db.Where(model).First(model)
	fmt.Println(model)

	return
}
