package models

import "github.com/jinzhu/gorm"

type TestModel struct {
	gorm.Model
	Code string
	Name string
}
