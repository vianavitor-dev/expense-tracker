package main

import (
	"github.com/spf13/cobra"
	"github.com/vianavitor-dev/expense-tracker/cmd"
)

func main() {
	var rootCmd = &cobra.Command{Use: "expense-tracker"}

	rootCmd.AddCommand(cmd.AddExpenseCommand())
	rootCmd.Execute()
}
