package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/templates"
	"github.com/cybera/ccds/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:              "init",
	Short:            "Creates a basic data science project skeleton",
	Args:             cobra.ExactArgs(0),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		licenses := []string{
			"MIT",
			"BSD-3-Clause",
			"None",
		}

		if viper.GetString("ProjectRoot") != "" {
			log.Fatal("Project has already been initialized")
		}

		projectRoot, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(os.Stdin)

		files, err := ioutil.ReadDir(projectRoot)
		if err != nil {
			log.Fatal(err)
		}

		if len(files) > 0 {
			fmt.Print("This directory is not empty, initialize anyways? [y/N]: ")

			for {
				input, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}

				input = strings.ToLower(utils.Chomp(input))

				if input == "y" {
					break
				} else if input == "n" || input == "" {
					os.Exit(0)
				}

				fmt.Print("Please answer [y/N]: ")
			}
		}

		viper.Set("ProjectRoot", projectRoot)

		// fmt.Print("Project name: ")
		// projectName, err := reader.ReadString('\n')
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// projectName = utils.Chomp(projectName)

		fmt.Print("Author (Your name or organization/company/team): ")
		author, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		author = utils.Chomp(author)

		var license, choices string

		for i := range licenses {
			choices += strconv.Itoa(i+1) + ", "
		}
		choices = choices[:len(choices)-2]

		fmt.Println("Select your license: ")
		for i, v := range licenses {
			fmt.Println(i+1, "-", v)
		}

		for {
			fmt.Printf("Choose %s: ", choices)
			choice, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			choice = utils.Chomp(choice)
			index, err := strconv.Atoi(choice)
			if err == nil && index > 0 && index <= len(licenses) {
				license = licenses[index-1]
				break
			}
		}

		if err := createSkeleton(); err != nil {
			log.Fatal(err)
		}

		if err := writeLicense(author, license); err != nil {
			log.Fatal(err)
		}

		if err := initRepo(); err != nil {
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

func createSkeleton() error {
	projectRoot := viper.GetString("ProjectRoot")
	language := viper.GetString("PrimaryLanguage")
	gitignore := "gitignore/" + language

	// Key is the directory path, value is whether to create a .gitkeep file
	directories := map[string]bool{
		".ccds":             true,
		"data":              false,
		"data/external":     true,
		"data/interim":      true,
		"data/processed":    true,
		"data/raw":          true,
		"docs":              true,
		"models":            true,
		"notebooks":         true,
		"references":        true,
		"reports":           false,
		"reports/figures":   true,
		"src":               false,
		"src/data":          true,
		"src/features":      true,
		"src/models":        true,
		"src/scripts":       true,
		"src/visualization": true,
	}

	files := map[string]string{
		gitignore:                   ".gitignore",
		"docker/Dockerfile":         paths.Dockerfile(projectRoot),
		"docker/docker-compose.yml": paths.DockerCompose(projectRoot),
	}

	for dir, keep := range directories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}

		if keep {
			path := filepath.Join(dir, ".gitkeep")
			file, err := os.Create(path)
			if err != nil {
				return errors.Wrapf(err, "failed to create file %s", path)
			}
			file.Close()
		}
	}

	for src, dest := range files {
		if err := templates.Write(src, dest, struct{}{}); err != nil {
			return err
		}

	}

	return nil
}

func writeLicense(author, license string) error {
	if license == "None" {
		return nil
	}

	src := "licenses/" + license

	data := struct {
		Year, Author string
	}{
		strconv.Itoa(time.Now().Year()),
		author,
	}

	return templates.Write(src, "LICENSE", data)
}

func initRepo() error {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return errors.Wrap(err, "failed to detect existing git repo")
	}

	for _, f := range files {
		if f.Name() == ".git" && f.IsDir() {
			return errors.New("git repo already exists")
		}
	}

	if _, err := exec.LookPath("git"); err != nil {
		return errors.Wrap(err, "git not found in path")
	}

	if err := exec.Command("git", "init").Run(); err != nil {
		return errors.Wrap(err, "failed to initialize git repo")
	}

	gitAdd(".ccds")
	gitCommit("Add ccds config directory")
	gitAdd(".gitignore", "LICENSE")
	gitCommit("Add standard repo files")
	gitAdd("Dockerfile", "docker-compose.yml")
	gitCommit("Add Docker configuration for Jupyter")
	gitAdd("data")
	gitCommit("Add directory for storing datasets")
	gitAdd("docs")
	gitCommit("Add directory for storing documentation")
	gitAdd("models")
	gitCommit("Add directory for storing models")
	gitAdd("notebooks")
	gitCommit("Add directory for storing notebooks")
	gitAdd("references")
	gitCommit("Add directory for storing references")
	gitAdd("reports")
	gitCommit("Add directory for storing reports")
	gitAdd("src")
	gitCommit("Add directory for storing source code")

	return nil
}

func gitAdd(paths ...string) error {
	args := append([]string{"add"}, paths...)
	return exec.Command("git", args...).Run()
}

func gitCommit(message string) error {
	return exec.Command("git", "commit", "-m", message).Run()
}
