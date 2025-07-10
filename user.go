package main

import (
	"github.com/rivo/tview"
	"golang.org/x/crypto/bcrypt"
)

// createUserLayout creates the layout for regular users
func createUserLayout(pages *tview.Pages, table *tview.Table, summary *tview.TextView) *tview.Flex {
	footer := createFooter(app, pages, table, summary)

	mainContent := tview.NewFlex().
		AddItem(table, 0, 1, true).
		AddItem(summary, 40, 0, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainContent, 0, 1, true).
		AddItem(footer, 1, 0, false)

	updateFinancialViews(table, summary)
	return layout
}

// createLoginForm creates the login form for users
func createLoginForm(pages *tview.Pages) *tview.Form {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil)

	form.AddButton("Login", func() {
		username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		login(pages, username, password)
	}).
		AddButton("Register", func() {
			regForm := createRegistrationForm(pages)
			pages.AddPage("register", regForm, true, true)
			pages.SwitchToPage("register")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Login").SetTitleAlign(tview.AlignLeft)
	return form
}

// createRegistrationForm creates the registration form for new users
func createRegistrationForm(pages *tview.Pages) *tview.Form {
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil)

	form.AddButton("Register", func() {
		username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := User{Name: username, Password: string(hashedPassword), Role: "user"}
		if err := db.Create(&user).Error; err != nil {
			// Handle error, e.g., show a modal
		} else {
			pages.SwitchToPage("login")
			pages.RemovePage("register")
		}
	}).
		AddButton("Back to Login", func() {
			pages.SwitchToPage("login")
			pages.RemovePage("register")
		})

	form.SetBorder(true).SetTitle("Register").SetTitleAlign(tview.AlignLeft)
	return form
} 