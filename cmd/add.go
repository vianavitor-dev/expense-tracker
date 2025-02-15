package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vianavitor-dev/expense-tracker/model"
)

var (
	amount      float64 = 0.0
	description string  = ""
)

func AddExpenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [-D description] [-A amount]",
		Short: "Add a new Expense on the file",
		Run: func(cmd *cobra.Command, args []string) {

			if amount == 0.0 && description == "" {
				fmt.Print("at least one field must be filled in")
				return
			}

			var currentDate = time.Now().Format("2006-01-02")

			var e = model.NewExpense(currentDate, description, amount)

			if err := e.Create(); err != nil {
				fmt.Print(err, "\t")
				return
			}

			fmt.Printf("expense added succesfully (ID: %d)\n", e.ID)
		},
	}

	cmd.Flags().Float64VarP(&amount, "amount", "A", 0.0, "amount of the expense")
	cmd.Flags().StringVarP(&description, "description", "D", "", "description of the expense")

	return cmd
}
