package main

import (
	"github.com/rivo/tview"
)

var adminSelectedUser *User

// createAdminLayout creates the layout for admin users with user list
func createAdminLayout(pages *tview.Pages, table *tview.Table, summary *tview.TextView) *tview.Flex {
	userList := createAdminUserList(pages, table, summary)
	
	mainContent := tview.NewFlex().
		AddItem(table, 0, 1, true).
		AddItem(summary, 40, 0, false)

	footer := createFooter(app, pages, table, summary)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainContent, 0, 1, true).
		AddItem(footer, 1, 0, false)

	adminLayout := tview.NewFlex().
		AddItem(userList, 0, 1, true).
		AddItem(layout, 0, 3, true)
	
	return adminLayout
}

// createAdminUserList creates the user list for admin interface
func createAdminUserList(pages *tview.Pages, table *tview.Table, summary *tview.TextView) *tview.List {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetTitle("Users").SetBorder(true)

	var users []User
	db.Find(&users)

	for i, user := range users {
		// Create a copy of the user to avoid closure issues
		userCopy := users[i]
		list.AddItem(user.Name, "", 0, func() {
			adminSelectedUser = &userCopy
			updateFinancialViews(table, summary)
		})
	}
	return list
}

// getSelectedUser returns the user whose transactions should be displayed
func getSelectedUser() User {
	if currentUser.Role == "admin" && adminSelectedUser != nil {
		return *adminSelectedUser
	}
	return currentUser
} 