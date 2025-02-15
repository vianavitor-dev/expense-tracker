package model_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/vianavitor-dev/expense-tracker/model"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func TestCreate(t *testing.T) {
	var e = model.NewExpense("", RandStringBytes(20), rand.Float64())
	currTime := time.Now()

	date := currTime.Format("2006-01-02")
	e.Date = date

	if err := e.Create(); err != nil {
		fmt.Print(err, "\t")
	} else {
		fmt.Printf("Expense added succesfully (ID: %d)\n", e.ID)
	}

}

func TestUpdate(t *testing.T) {

	var e = &model.Expense{
		Date:        "2025-02-15",
		Description: "uga",
	}

	var id = int64(1)

	if err := e.Update(id); err != nil {
		fmt.Print(err, "\t")
	} else {
		fmt.Printf("Expense modified successfully (ID: %d)\n", id)
	}
}

func TestDelete(t *testing.T) {
	var e *model.Expense
	var id = int64(4)

	if err := e.Delete(id); err != nil {
		fmt.Print(err, "\t")
	} else {
		fmt.Print("Expense deleted successfully\n")
	}
}

func TestViewAll(t *testing.T) {
	var e *model.Expense

	if err := e.ViewAll(); err != nil {
		fmt.Print(err, "\t")
	}
}

func TestViewSumary(t *testing.T) {
	var e *model.Expense

	if err := e.ViewSumary(0); err != nil {
		fmt.Print(err, "\t")
	}
}
