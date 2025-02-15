package model

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type Expense struct {
	ID          int64   `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

var pathName = "expense-list"

func NewExpense(date string, desc string, amount float64) *Expense {
	e := new(Expense)
	e.Date = date
	e.Description = desc
	e.Amount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)

	return e
}

func (e *Expense) Create() error {

	file, err := os.OpenFile(pathName+".json", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	var expenses []Expense
	decoder := json.NewDecoder(file)

	decoder.Token()

	for decoder.More() {
		var currentExp Expense

		decoder.Decode(&currentExp)
		expenses = append(expenses, currentExp)
	}

	e.ID = int64(len(expenses) + 1)
	expenses = append(expenses, *e)

	file, err = os.Create(pathName + ".json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(expenses)
	if err != nil {
		return err
	}

	return nil
}

func (e *Expense) Update(id int64) error {

	file, err := os.OpenFile(pathName+".json", os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fInfo, _ := os.Stat(pathName + ".json")

	bs := make([]byte, fInfo.Size())

	_, err = bufio.NewReader(file).Read(bs)
	if err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	file.Seek(0, 0)

	var expenses []Expense
	if err := json.Unmarshal(bs, &expenses); err != nil {

		if len(expenses) == 0 {
			return fmt.Errorf("file is empty")
		}

		return err
	}

	var found = false
	for i := range expenses {
		if expenses[i].ID == id {

			if e.Description != "" {
				expenses[i].Description = e.Description
			}
			if e.Date != "" {
				expenses[i].Date = e.Date
			}
			if e.Amount > 0 {
				expenses[i].Amount = e.Amount
			}

			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expense not found")
	}

	if err = json.NewEncoder(file).Encode(expenses); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (e *Expense) Delete(id int64) error {

	file, err := os.Open(pathName + ".json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var expenses []Expense
	if err = json.NewDecoder(file).Decode(&expenses); err != nil {

		if len(expenses) == 0 {
			return fmt.Errorf("file is empty")
		}

		return err
	}

	newExpenses := make([]Expense, len(expenses)-1)
	index := 0
	found := false

	for i := range expenses {
		if expenses[i].ID != id {
			newExpenses[index] = expenses[i]

			if index >= len(newExpenses)-1 {
				break
			}

			index++
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("expense not found")
	}

	if file, err = os.Create(pathName + ".json"); err != nil {
		log.Fatal(err)
	}

	if len(newExpenses) != 0 {
		if err = json.NewEncoder(file).Encode(newExpenses); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (e *Expense) ViewAll() error {

	file, err := os.Open(pathName + ".json")
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file not found")
		}

		log.Fatal(err)
	}

	defer file.Close()

	var decoder = json.NewDecoder(file)

	if !decoder.More() {
		return fmt.Errorf("file is empty")
	}

	decoder.Token()

	var currentExp Expense
	tabW := new(tabwriter.Writer)

	tabW.Init(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprint(tabW, "ID\tDate\tDescription\tAmount\n")

	for decoder.More() {
		if err := decoder.Decode(&currentExp); err != nil {
			return err
		}

		fmt.Fprintf(tabW, "%d\t%s\t%s\t%v \n", currentExp.ID, currentExp.Date, currentExp.Description, currentExp.Amount)
	}

	fmt.Fprintln(tabW)

	tabW.Flush()
	return nil
}

func (e *Expense) ViewSumary(month uint) error {

	file, err := os.Open(pathName + ".json")
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file not found")
		}

		log.Fatal(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var expense Expense
	var totalExpenses float64

	decoder.Token()

	if !decoder.More() {
		return fmt.Errorf("file is empty")
	}

	for decoder.More() {
		if err := decoder.Decode(&expense); err != nil {
			log.Fatal(err)
		}

		if month > 0 {
			var currentTime = time.Now()
			var expenseDate, err = time.Parse("2006-01-02", expense.Date)
			if err != nil {
				return err
			}

			if expenseDate.Month() == time.Month(month) && expenseDate.Year() == currentTime.Year() {
				totalExpenses += expense.Amount
			}

		} else {
			totalExpenses += expense.Amount
		}
	}

	fmt.Printf("Total expenses: $%.2f", totalExpenses)
	return nil
}
