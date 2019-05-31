package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cybera/ccds/internal/paths"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ccds",
	Short: "CCDS is a data science project generator",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		projectRoot, err := paths.ProjectRoot()
		if err != nil {
			log.Fatal("Not under a valid project directory")
		}

		viper.SetDefault("datasets", map[string]interface{}{})

		viper.SetConfigFile(filepath.Join(projectRoot, paths.ProjectMetadata()))
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("failed to read config: ", err)
		}

		viper.Set("ProjectRoot", projectRoot)
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
