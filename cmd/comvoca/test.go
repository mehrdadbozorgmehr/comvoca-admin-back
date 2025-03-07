package cmd

import (
	"fmt"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/spf13/cobra"
	"os/exec"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running tests...")
		// Add your code to run tests here
		command := exec.Command("make", "test")
		command.Stdout = cmd.OutOrStdout()
		command.Stderr = cmd.OutOrStderr()
		if err := command.Run(); err != nil {
			logger.Error("Error running tests:", err)
		}
	},
}
