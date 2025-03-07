package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/spf13/cobra"
)

var InsertTestDataCmd = &cobra.Command{
	Use:   "insert-testdata",
	Short: "Insert test data into the database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Inserting test data...")
		// Add your code to insert test data here
		command := exec.Command("make", "testdata")
		command.Stdout = cmd.OutOrStdout()
		command.Stderr = cmd.OutOrStderr()
		if err := command.Run(); err != nil {
			logger.Error("Error inserting test data:", err)
		}
	},
}
