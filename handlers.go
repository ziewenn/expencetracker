package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/crypto/bcrypt"
)

func login(pages *tview.Pages, username, password string) {
	go func() {
		var user User
		if err := db.Where("name = ?", username).First(&user).Error; err != nil {
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

			return
		}

		app.QueueUpdateDraw(func() {
			currentUser = user
			appLayout := createAppLayout(pages)
			pages.AddPage("app", appLayout, true, true)
			pages.SwitchToPage("app")
		})
	}()
}

func updateFinancialViews(table *tview.Table, summary *tview.TextView) {
	go func() {
		user := getSelectedUser()

		var transactions []Transaction
		db.Where("user_id = ?", user.ID).Order("created_at desc").Find(&transactions)

		var totalIncome, totalExpenses, maxExpense float64
		var maxExpenseCategory string = "N/A"

		var maxExpenseTransaction Transaction
		if err := db.Where("user_id = ? AND type = ?", user.ID, "expense").Order("amount desc").First(&maxExpenseTransaction).Error; err == nil {
			maxExpense = maxExpenseTransaction.Amount
			maxExpenseCategory = maxExpenseTransaction.Category
		}

		db.Model(&Transaction{}).Where("user_id = ? AND type = ?", user.ID, "income").Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalIncome)
		db.Model(&Transaction{}).Where("user_id = ? AND type = ?", user.ID, "expense").Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalExpenses)

		balance := totalIncome - totalExpenses

		app.QueueUpdateDraw(func() {

			table.Clear()
			headers := []string{"ID", "Amount", "Category", "Type", "Description", "Date"}
			for i, header := range headers {
				table.SetCell(0, i, tview.NewTableCell(header).
					SetTextColor(tcell.ColorYellow).
					SetAlign(tview.AlignCenter))
			}
			for i, t := range transactions {
				table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", t.ID)))
				table.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%.2f", t.Amount)))
				table.SetCell(i+1, 2, tview.NewTableCell(t.Category))
				table.SetCell(i+1, 3, tview.NewTableCell(t.Type))
				table.SetCell(i+1, 4, tview.NewTableCell(t.Description))
				table.SetCell(i+1, 5, tview.NewTableCell(t.CreatedAt.Format("2006-01-02")))
			}

			var summaryText strings.Builder
			summaryText.WriteString(fmt.Sprintf("[green]Total Income:   [white]%.2f\n", totalIncome))
			summaryText.WriteString(fmt.Sprintf("[red]Total Expenses: [white]%.2f\n", totalExpenses))
			summaryText.WriteString(fmt.Sprintf("[blue]Max Expense: [white]%.2f\n", maxExpense))
			summaryText.WriteString(fmt.Sprintf("[blue]Category: [white]%s\n", maxExpenseCategory))
			summaryText.WriteString("\n--------------------\n")
			summaryText.WriteString(fmt.Sprintf("\n[yellow]Balance:        [white]%.2f\n", balance))

			barWidth := 30
			summaryText.WriteString("\n\n")
			summaryText.WriteString(fmt.Sprintf("[green]Income: [white]%.2f\n", totalIncome) + "[green]" + generateBar(totalIncome, totalIncome, barWidth) + "[white]\n")
			summaryText.WriteString(fmt.Sprintf("[red]Expenses: [white]%.2f\n", totalExpenses) + "[red]" + generateBar(totalExpenses, totalIncome, barWidth) + "[white]\n")
			summaryText.WriteString(fmt.Sprintf("[yellow]Balance:  [white]%.2f\n", balance) + "[yellow]" + generateBar(balance, totalIncome, barWidth) + "[white]\n")

			summary.SetText(summaryText.String())
		})
	}()
}

func saveTransaction(form *tview.Form, pages *tview.Pages, table *tview.Table, summary *tview.TextView, transaction *Transaction) {
	amountStr := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		modal := tview.NewModal().
			SetText("Error: Please enter a valid positive amount").
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("error_modal")
			})
		pages.AddPage("error_modal", modal, true, true)
		pages.SwitchToPage("error_modal")
		return
	}
	
	category := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
	if category == "" {
		modal := tview.NewModal().
			SetText("Error: Please enter a category").
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("error_modal")
			})
		pages.AddPage("error_modal", modal, true, true)
		pages.SwitchToPage("error_modal")
		return
	}
	
	_, tType := form.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()
	description := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()

	if tType == "expense" {
		user := getSelectedUser()

		var totalIncome float64
		db.Model(&Transaction{}).Where("user_id = ? AND type = ?", user.ID, "income").Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalIncome)

		var totalExpenses float64
		expenseQuery := db.Model(&Transaction{}).Where("user_id = ? AND type = ?", user.ID, "expense")
		if transaction != nil {
			expenseQuery = expenseQuery.Where("id != ?", transaction.ID)
		}
		expenseQuery.Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalExpenses)

		if totalExpenses+amount > totalIncome {
			modal := tview.NewModal().
				SetText("Error: Total expenses cannot exceed total income.").
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
									pages.RemovePage("error_modal")
			})
		app.QueueUpdateDraw(func() {
			pages.AddPage("error_modal", modal, true, true)
			pages.SwitchToPage("error_modal")
		})
			return
		}
	}

	go func() {
		if transaction == nil {

			user := getSelectedUser()
			newTransaction := Transaction{
				Amount:      amount,
				Category:    category,
				Type:        tType,
				Description: description,
				UserID:      user.ID,
			}
			db.Create(&newTransaction)
		} else {

			transaction.Amount = amount
			transaction.Category = category
			transaction.Type = tType
			transaction.Description = description
			db.Save(transaction)
		}
		updateFinancialViews(table, summary)
	}()

	pages.SwitchToPage("app")
	pages.RemovePage("form")
}

func editTransaction(pages *tview.Pages, table *tview.Table, summary *tview.TextView) {
	row, _ := table.GetSelection()
	if row == 0 { 
		return
	}
	
	// Check if we have a valid cell
	cell := table.GetCell(row, 0)
	if cell == nil {
		return
	}
	
	idStr := cell.Text
	if idStr == "" {
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Show error modal
		modal := tview.NewModal().
			SetText("Error: Invalid transaction ID").
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("error_modal")
			})
		pages.AddPage("error_modal", modal, true, true)
		pages.SwitchToPage("error_modal")
		return 
	}

	var transaction Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		// Show error modal
		modal := tview.NewModal().
			SetText("Error: Transaction not found").
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("error_modal")
			})
		pages.AddPage("error_modal", modal, true, true)
		pages.SwitchToPage("error_modal")
		return 
	}

	form := createForm(pages, table, summary, &transaction)
	pages.AddPage("form", form, true, true)
	pages.SwitchToPage("form")
}

func showDeleteConfirmation(pages *tview.Pages, table *tview.Table, summary *tview.TextView) {
	row, _ := table.GetSelection()
	if row == 0 {
		return
	}
	
	// Check if we have a valid cell
	cell := table.GetCell(row, 0)
	if cell == nil {
		return
	}
	
	idStr := cell.Text
	if idStr == "" {
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete transaction ID %d?", id)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				go func() {
					db.Delete(&Transaction{}, id)
					updateFinancialViews(table, summary)
				}()
			}
			pages.SwitchToPage("app")
			pages.RemovePage("modal")
		})

	pages.AddPage("modal", modal, true, true)
	pages.SwitchToPage("modal")
} 