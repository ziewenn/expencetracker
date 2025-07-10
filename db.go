package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Password     string
	Role         string
	Transactions []Transaction
}

type Transaction struct {
	gorm.Model
	Amount      float64
	Category    string
	Description string
	Type        string
	UserID      uint
}

var db *gorm.DB
var currentUser User

func (t *Transaction) GetCategory() string {
	if t == nil {
		return ""
	}
	return t.Category
}

func (t *Transaction) GetDescription() string {
	if t == nil {
		return ""
	}
	return t.Description
}

func (t *Transaction) GetTypeIndex() int {
	if t == nil || t.Type == "income" {
		return 0
	}
	return 1
} 