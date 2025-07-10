package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createAppLayout(pages *tview.Pages) *tview.Flex {
	table := createTable()
	summary := createSummaryView()

	if currentUser.Role == "admin" {
		return createAdminLayout(pages, table, summary)
	}

	return createUserLayout(pages, table, summary)
}

func createTable() *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table.SetTitle("Transactions").SetBorder(true)

	headers := []string{"ID", "Amount", "Category", "Type", "Description", "Date"}
	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))
	}

	return table
}

func createFooter(app *tview.Application, pages *tview.Pages, table *tview.Table, summary *tview.TextView) *tview.Flex {
	newButton := tview.NewButton("New").SetSelectedFunc(func() {
		form := createForm(pages, table, summary, nil)
		pages.AddPage("form", form, true, true)
		pages.SwitchToPage("form")
	})

	editButton := tview.NewButton("Edit").SetSelectedFunc(func() {
		editTransaction(pages, table, summary)
	})

	deleteButton := tview.NewButton("Delete").SetSelectedFunc(func() {
		showDeleteConfirmation(pages, table, summary)
	})

	logoutButton := tview.NewButton("Logout").SetSelectedFunc(func() {
		// Clear current user and admin selection
		currentUser = User{}
		adminSelectedUser = nil
		
		// Remove app page and return to login
		pages.RemovePage("app")
		pages.SwitchToPage("login")
	})

	quitButton := tview.NewButton("Quit").SetSelectedFunc(func() {
		app.Stop()
	})

	footer := tview.NewFlex().
		AddItem(newButton, 0, 1, false).
		AddItem(editButton, 0, 1, false).
		AddItem(deleteButton, 0, 1, false).
		AddItem(logoutButton, 0, 1, false).
		AddItem(quitButton, 0, 1, false)

	return footer
}

func createSummaryView() *tview.TextView {
	summary := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	summary.SetBorder(true).SetTitle("Financial Summary")
	return summary
}

func generateBar(value, maxValue float64, width int) string {
	if maxValue == 0 {
		return strings.Repeat(" ", width)
	}
	ratio := value / maxValue
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	blocks := int(ratio * float64(width))
	return strings.Repeat("█", blocks) + strings.Repeat(" ", width-blocks)
}

func createForm(pages *tview.Pages, table *tview.Table, summary *tview.TextView, transaction *Transaction) *tview.Form {
	form := tview.NewForm()

	amountStr := ""
	if transaction != nil {
		amountStr = fmt.Sprintf("%.2f", transaction.Amount)
	}

	form.AddInputField("Amount", amountStr, 20, tview.InputFieldFloat, nil).
		AddInputField("Category", transaction.GetCategory(), 20, nil, nil).
		AddDropDown("Type", []string{"income", "expense"}, transaction.GetTypeIndex(), nil).
		AddInputField("Description", transaction.GetDescription(), 40, nil, nil)

	if transaction == nil {
		form.AddButton("Save", func() {
			saveTransaction(form, pages, table, summary, nil)
		})
	} else {
		form.AddButton("Update", func() {
			saveTransaction(form, pages, table, summary, transaction)
		})
	}

	form.AddButton("Cancel", func() {
		pages.SwitchToPage("app")
		pages.RemovePage("form")
	})

	form.SetBorder(true)
	if transaction == nil {
		form.SetTitle("Add New Transaction")
	} else {
		form.SetTitle("Edit Transaction")
	}

	return form
} 