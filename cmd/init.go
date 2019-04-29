package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a basic data science project skeleton",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		licenses := map[string]string{
			"1": "MIT",
			"2": "BSD-3-Clause",
			"3": "None",
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Project name: ")
		projectName, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		projectName = chomp(projectName)

		fmt.Print("Author (Your name or organization/company/team): ")
		author, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		author = chomp(author)

		fmt.Println("Select your license: ")
		for k, v := range licenses {
			fmt.Println(k, "-", v)
		}
		fmt.Print("Choose from 1, 2, 3: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = chomp(choice)
		license, ok := licenses[choice]
		if !ok {
			panic("Bad choice!")
		}

		createSkeleton()

		licenseText, err := ioutil.ReadFile("templates/licenses/" + license)
		if err != nil {
			panic(err)
		}
		tmpl, err := template.New("License").Parse(string(licenseText))
		if err != nil {
			panic(err)
		}

		data := struct {
			Year, Author string
		}{
			strconv.Itoa(time.Now().Year()),
			author,
		}

		licenseFile, err := os.Create("foo/LICENSE")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(licenseFile, data)
		if err != nil {
			panic(err)
		}
		licenseFile.Close()

		fmt.Println(projectName, author, license)
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
		"foo/.ccds",
		"foo/data/external",
		"foo/data/interim",
		"foo/data/processed",
		"foo/data/raw",
		"foo/docs",
		"foo/models",
		"foo/notebooks",
		"foo/references",
		"foo/reports/figures",
		"foo/src/data",
		"foo/src/features",
		"foo/src/models",
		"foo/src/visualization",
	}

	files := map[string]string{
		"templates/docker/Dockerfile":         "foo/Dockerfile",
		"templates/docker/docker-compose.yml": "foo/docker-compose.yml",
	}

	for _, dir := range directories {
		os.MkdirAll(dir, os.ModePerm)
	}

	for src, dest := range files {
		input, err := os.Open(src)
		if err != nil {
			panic(err)
		}

		output, err := os.Create(dest)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(output, input)
		if err != nil {
			panic(err)
		}
	}
}

func chomp(s string) string {
	return strings.Trim(s, " \r\n")
}
