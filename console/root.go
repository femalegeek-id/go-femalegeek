package console

import (
	"fmt"
	"os"

	"femalegeek/config"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cobra-example",
	Short: "An example of cobra",
	Long: `This application shows how to create modern CLI 
			applications in go using Cobra CLI library`,
}

// Execute RootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.GetConf()
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &log.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)

	logLevel, err := log.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

}
