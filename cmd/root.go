package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cybera/mimir/internal/paths"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nonInteractive, yesToAll bool

var rootCmd = &cobra.Command{
	Use:   "mimir",
	Short: "Mimir is a data science project generator",
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

func init() {
	rootCmd.PersistentFlags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "Error if any user input is required")
	rootCmd.PersistentFlags().BoolVarP(&yesToAll, "yes", "y", false, "Answer yes to any prompts automatically")
}
