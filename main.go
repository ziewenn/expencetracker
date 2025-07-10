package main

import (
	"log"

	"github.com/rivo/tview"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var app *tview.Application

func main() {
	dsn := "DATABASE_DSN"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&User{}, &Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	createAdminIfNotExists()

	app = tview.NewApplication()
	pages := tview.NewPages()

	loginForm := createLoginForm(pages)
	pages.AddPage("login", loginForm, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func createAdminIfNotExists() {
	var admin User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
			admin = User{Name: "admin", Password: string(hashedPassword), Role: "admin"}
			db.Create(&admin)
		}
	}
}