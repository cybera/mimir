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
	"time"

	"github.com/cybera/mimir/internal/languages"
	"github.com/cybera/mimir/internal/paths"
	"github.com/cybera/mimir/internal/templates"
	"github.com/cybera/mimir/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var author, license, language string
var force bool

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

		if len(files) > 0 && !force {
			if !utils.GetYesNo(reader, "This directory is not empty, initialize anyways?", false, nonInteractive) {
				os.Exit(0)
			}
		}

		viper.Set("ProjectRoot", projectRoot)

		if author == "" {
			fmt.Print("Author (Your name or organization/company/team): ")
			author = utils.GetInput(reader, nonInteractive)
		}

		if license == "" {
			license = ask(reader, "Select your license: ", licenses, -1)
		} else if !utils.Contains(licenses, license) {
			log.Fatal("unknown license")
		}

		if language == "" {
			language = ask(reader, "Select your language: ", languages.Supported, 0)
		} else if !utils.Contains(languages.Supported, language) {
			log.Fatal("unknown language")
		}

		viper.Set("Author", author)
		viper.Set("License", license)
		viper.Set("PrimaryLanguage", language)

		if err := initProject(projectRoot, author, license, language); err != nil {
			log.Fatalf("project initialization failed: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&author, "author", "", "Author name")
	initCmd.Flags().StringVar(&license, "license", "", "Project license")
	initCmd.Flags().StringVar(&language, "language", "", "Which programming language to use")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Ignore existing files and directories")
}

func ask(reader *bufio.Reader, text string, choices []string, def int) string {
	var numbers string
	for i := range choices {
		numbers += strconv.Itoa(i+1) + ", "
	}
	numbers = numbers[:len(numbers)-2]

	fmt.Println(text)
	for i, v := range choices {
		fmt.Println(i+1, "-", v)
	}

	var choice int

	for {
		if def >= 0 {
			def++
			fmt.Printf("Choose %s [%d]: ", numbers, def)
		} else {
			fmt.Printf("Choose %s: ", numbers)
		}
		input := utils.GetInput(reader, nonInteractive)

		if def > 0 && input == "" {
			choice = def
			break
		}

		var err error

		choice, err = strconv.Atoi(input)
		if err == nil && choice > 0 && choice <= len(choices) {
			break
		}
	}

	return choices[choice-1]
}

func getInput(reader *bufio.Reader) string {
	if nonInteractive {
		log.Fatal("\nerror: input required in non-interactive mode")
	}

	input, _ := reader.ReadString('\n')
	return utils.Chomp(input)
}

func initProject(projectRoot, author, license, language string) error {
	log.Println("Creating project skeleton...")
	if err := createSkeleton(projectRoot, language); err != nil {
		return err
	}

	if err := writeLicense(author, license); err != nil {
		return err
	}

	log.Println("Initializing git repository...")
	if err := initRepo(); err != nil {
		return err
	}

	return nil
}

func createSkeleton(projectRoot, language string) error {
	// Key is the directory path, value is whether to create a .gitkeep file
	directories := map[string]bool{
		".mimir":             false,
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
		"src/datasets":      true,
		"src/features":      true,
		"src/models":        true,
		"src/scripts":       true,
		"src/visualization": true,
	}

	files := map[string]string{
		"docker/Dockerfile":         filepath.Join(projectRoot, paths.Dockerfile()),
		"docker/docker-compose.yml": filepath.Join(projectRoot, paths.DockerCompose()),
		"project-settings.toml":     filepath.Join(projectRoot, paths.ExampleProjectSettings()),
	}

	for k, v := range languages.InitFiles[language] {
		files[k] = v
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

	data := struct {
		ProjectSettingsPath string
	}{
		ProjectSettingsPath: filepath.Join("../", paths.ProjectSettings()),
	}

	for src, dest := range files {
		if err := templates.WriteFile(src, dest, data); err != nil {
			return err
		}
	}

	if err := writeGitignore(language); err != nil {
		return err
	}

	if err := utils.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func writeGitignore(language string) error {
	sources := []string{"gitignore/general", "gitignore/" + language}

	file, err := os.Create(".gitignore")
	if err != nil {
		return errors.Wrapf(err, "failed to create file .gitignore")
	}
	defer file.Close()

	for _, s := range sources {
		if err := templates.Write(s, file, struct{}{}); err != nil {
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

	return templates.WriteFile(src, "LICENSE", data)
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

	gitAdd(".mimir")
	gitCommit("Add mimir config directory")
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
