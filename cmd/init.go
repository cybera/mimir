package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a basic data science project skeleton",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		licenses := map[string]string{
			"1": "MIT",
			"2": "BSD-3-Clause",
			"3": "None",
		}

		if _, err := paths.ProjectRootSafe(); err == nil {
			log.Fatal("Project has already been initialized")
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Project name: ")
		projectName, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		projectName = utils.Chomp(projectName)

		fmt.Print("Author (Your name or organization/company/team): ")
		author, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		author = utils.Chomp(author)

		fmt.Println("Select your license: ")
		for k, v := range licenses {
			fmt.Println(k, "-", v)
		}
		fmt.Print("Choose from 1, 2, 3: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		choice = utils.Chomp(choice)
		license, ok := licenses[choice]
		if !ok {
			log.Fatal(fmt.Sprintf("%s is not a valid choice!", choice))
		}

		createSkeleton()

		licenseText, err := templates.FindString("licenses/" + license)
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Year, Author string
		}{
			strconv.Itoa(time.Now().Year()),
			author,
		}

		tmpl, err := template.New("License").Parse(licenseText)
		if err != nil {
			log.Fatal(err)
		}

		licenseFile, err := os.Create("LICENSE")
		if err != nil {
			log.Fatal(err)
		}
		defer licenseFile.Close()

		err = tmpl.Execute(licenseFile, data)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createSkeleton() {
	directories := []string{
		".ccds",
		"data/external",
		"data/interim",
		"data/processed",
		"data/raw",
		"docs",
		"models",
		"notebooks",
		"references",
		"reports/figures",
		"src/data",
		"src/features",
		"src/models",
		"src/scripts",
		"src/visualization",
	}

	files := map[string]string{
		"docker/Dockerfile":         "Dockerfile",
		"docker/docker-compose.yml": "docker-compose.yml",
	}

	for _, dir := range directories {
		os.MkdirAll(dir, os.ModePerm)
	}

	for src, dest := range files {
		contents, err := templates.Find(src)
		if err != nil {
			log.Fatal(err)
		}

		output, err := os.Create(dest)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()

		_, err = output.Write(contents)
		if err != nil {
			log.Fatal(err)
		}
	}
}
