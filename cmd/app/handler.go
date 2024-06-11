package main

import (
	// "expense/pkg/models/mysql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/login.page.tmpl",
		"ui/html/base.layout.tmpl",
	}
	app.render(w, files, nil)

}
func (app *Application) AddUser(w http.ResponseWriter, r *http.Request) {

	name := "Aswin"
	email := "aswinbaiju4"
	password := "123"

	err := app.Users.Insert(name, email, password)
	if err != nil {
		app.ErrorLog.Fatal()
		return
	} else {
		app.Session.Put(r, "flash", "Task successfully created!")
	}

}

func (app *Application) AddSplit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	amount := r.FormValue("amount")
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		app.Session.Put(r, "flash", "Invalid amount!")
		return
	}
	note := r.FormValue("note")

	result, err := app.Expense.Insert(note, amountFloat, app.Session.GetInt(r, "userId"))
	http.Redirect(w, r, "/submit_expense", http.StatusSeeOther)
	if err != nil {
		app.ErrorLog.Fatal()
		return
	}
	app.Session.Put(r, "flash", "Task successfully created!")

	usersSelected := r.Form["user[]"]

	fmt.Println("Ids selected:")
	for _, id := range usersSelected {
		fmt.Println(id)
	}
	expenseId, err := result.LastInsertId()
	if err != nil {
		app.ErrorLog.Fatal()
	}
	log.Println("done.....")
	app.Expense.Insert2Split(expenseId, amountFloat, usersSelected)

}

func (app *Application) GetAddSplitForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/split.page.tmpl",
		"ui/html/base.layout.tmpl",
	}
	userList, errGettingList := app.Users.GetAllUsers()
	if errGettingList != nil {
		app.ErrorLog.Fatal()
		return
	}
	app.render(w, files, &templateData{
		UserData: userList,
	})

}
