package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cybera/ccds/internal/paths"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templates *packr.Box

var rootCmd = &cobra.Command{
	Use:   "ccds",
	Short: "CCDS is a data science project generator",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString("ProjectRoot") == "" {
			log.Fatal("Not under a valid project directory")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initPackr)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ccds.yaml)")
}

func initConfig() {
	projectRoot, err := paths.ProjectRoot()
	if err != nil {
		projectRoot = ""
	}

	viper.Set("ProjectRoot", projectRoot)
	viper.Set("ContainerRoot", paths.ContainerRoot())
}

func initPackr() {
	templates = packr.New("Templates", "../templates")
}
