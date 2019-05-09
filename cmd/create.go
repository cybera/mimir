package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cybera/ccds/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:              "create [author] [license]",
	Short:            "Same as the init command, but with arguments",
	Args:             cobra.ExactArgs(2),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		licenses := map[string]string{
			"1": "MIT",
			"2": "BSD-3-Clause",
			"3": "None",
		}

		if viper.GetString("ProjectRoot") != "" {
			log.Fatal("Project has already been initialized")
		}

		projectRoot, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		viper.Set("ProjectRoot", projectRoot)

		// fmt.Print("Project name: ")
		// projectName, err := reader.ReadString('\n')
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// projectName = utils.Chomp(projectName)

		author := utils.Chomp(args[0])
		choice := utils.Chomp(args[1])
		license, ok := licenses[choice]

		if !ok {
			log.Fatal(fmt.Sprintf("%s is not a valid choice!", choice))
		}

		createSkeleton(author, license)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
