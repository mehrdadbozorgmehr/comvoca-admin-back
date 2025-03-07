package main

import (
	cmd "github.com/Comvoca-AI/comvoca-admin-back/cmd/comvoca"
	"github.com/Comvoca-AI/comvoca-admin-back/config"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "comvoca-admin-back",
	Short: "Comvoca Admin Back",
	Long:  `	This is a back-end server for Comvoca Admin.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	logger.InitLogger()
	config.InitConfig()
	rootCmd.AddCommand(cmd.RunCmd)
	rootCmd.AddCommand(cmd.TestCmd)
	rootCmd.AddCommand(cmd.InsertTestDataCmd)
}

// @title Comvoca Admin API
// @version 1.0
// @description This is a back-end server for Comvoca Admin.
// @termsOfService http://swagger.io/terms/
// @contact.name info@comvoca.com
// @contact.email info@comvoca.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	if config.AppConfig.Server.MemoryProfile == true {
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	} else if config.AppConfig.Server.CPUProfile == true {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	Execute()
}
